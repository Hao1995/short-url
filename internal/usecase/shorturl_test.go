package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"testing"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/mocks/internal_/usecase"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"github.com/viney-shih/go-cache"
)

type ShortUrlUseCaseTestSuite struct {
	suite.Suite
	ctx               context.Context
	now               time.Time
	dockertestClose   func() error
	cacheFactoryClose func()
	host              string
	port              string

	ring *redis.Ring
	repo *usecase.Repository
	impl UseCase
}

func TestShortUrlUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlUseCaseTestSuite))
}

func (s *ShortUrlUseCaseTestSuite) SetupSuite() {
	loc, _ := time.LoadLocation("")
	time.Local = loc

	s.now = time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return s.now
	}

	// Run Redis
	var err error
	s.host, s.port, s.dockertestClose, err = ConnectToDockerTestRedis()
	if err != nil {
		log.Fatal("failed to set up redis container: ", err)
	}
}

func (s *ShortUrlUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ShortUrlUseCaseTestSuite) SetupSubTest() {
	// Reset after each sub-test in order to rest local cache

	// Setup Cache
	tinyLfu := cache.NewTinyLFU(10000)

	s.ring = redis.NewRing(&redis.RingOptions{Addrs: map[string]string{s.host: ":" + s.port}})
	rds := cache.NewRedis(s.ring)

	cacheFactory := cache.NewFactory(rds, tinyLfu)
	s.cacheFactoryClose = cacheFactory.Close

	cacheIns := cacheFactory.NewCache([]cache.Setting{
		{
			Prefix: domain.CACHE_PREFIX_SHORT_URL,
			CacheAttributes: map[cache.Type]cache.Attribute{
				cache.SharedCacheType: {TTL: time.Hour},
				cache.LocalCacheType:  {TTL: 10 * time.Second},
			},
			MarshalFunc:   json.Marshal,
			UnmarshalFunc: json.Unmarshal,
		},
	})

	s.repo = usecase.NewRepository(s.T())
	s.impl = NewShortUrlUseCase(s.repo, cacheIns)
}

func (s *ShortUrlUseCaseTestSuite) TearDownSubTest() {
	// clear registered prefix
	cache.ClearPrefix()

	// clean up all in redis
	s.Require().NoError(s.ring.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
		return client.FlushDB(ctx).Err()
	}))

	s.cacheFactoryClose()
}

func (s *ShortUrlUseCaseTestSuite) TearDownTest() {}

func (s *ShortUrlUseCaseTestSuite) TearDownSuite() {
	s.dockertestClose()
}

func (s *ShortUrlUseCaseTestSuite) TestCreate() {
	for _, t := range []struct {
		name   string
		req    *domain.CreateReqDto
		setup  func()
		exp    *domain.CreateRespDto
		expErr error
	}{
		{
			name: "create a record successfully",
			req: &domain.CreateReqDto{
				Url:      "https://example.com/whatever1",
				ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
			},
			setup: func() {
				targetID := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte("https://example.com/whatever1")))
				s.repo.On("Create", s.ctx, &domain.CreateReqDto{
					Url:      "https://example.com/whatever1",
					TargetID: targetID,
					ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
				}).Once().Return("testid1", nil)
			},
			exp: &domain.CreateRespDto{
				TargetID: "testid1",
				ShortUrl: "http://localhost/testid1",
			},
			expErr: nil,
		},
		{
			name: "create a duplicated record successfully",
			req: &domain.CreateReqDto{
				Url:      "https://example.com/whatever1",
				ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
			},
			setup: func() {
				// the first record
				targetID := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte("https://example.com/whatever1")))
				s.repo.On("Create", s.ctx, &domain.CreateReqDto{
					Url:      "https://example.com/whatever1",
					TargetID: targetID,
					ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
				}).Once().Return("", domain.ErrDuplicatedKey)

				// added suffix record
				suffix := "whatever"
				randString = func() string {
					return suffix
				}
				targetID = fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte("https://example.com/whatever1"+suffix)))
				s.repo.On("Create", s.ctx, &domain.CreateReqDto{
					Url:      "https://example.com/whatever1",
					TargetID: targetID,
					ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
				}).Once().Return("testid2", nil)
			},
			exp: &domain.CreateRespDto{
				TargetID: "testid2",
				ShortUrl: "http://localhost/testid2",
			},
			expErr: nil,
		},
		{
			name: "failed to create a record due to unknown error",
			req: &domain.CreateReqDto{
				Url:      "https://example.com/whatever2",
				TargetID: "testid1",
				ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
			},
			setup: func() {
				targetID := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte("https://example.com/whatever2")))
				s.repo.On("Create", s.ctx, &domain.CreateReqDto{
					Url:      "https://example.com/whatever2",
					TargetID: targetID,
					ExpireAt: time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC),
				}).Once().Return("", errors.New("unknown error"))
			},
			exp:    nil,
			expErr: errors.New("unknown error"),
		},
	} {
		s.Suite.Run(t.name, func() {
			ctx := context.Background()
			if t.setup != nil {
				t.setup()
			}
			id, err := s.impl.Create(ctx, t.req)
			s.Equal(err, t.expErr)
			s.Equal(t.exp, id)
		})
	}
}

func (s *ShortUrlUseCaseTestSuite) TestGet() {
	for _, t := range []struct {
		name   string
		req    string
		setup  func()
		check  func()
		expObj *domain.GetRespDto
		expErr error
	}{
		{
			name: "get record successfully when the record exist",
			req:  "testid1",
			setup: func() {
				s.repo.On("Get", s.ctx, "testid1").Once().Return(&domain.GetRespDto{
					Url:      "https://example.com/whatever1",
					ExpireAt: s.now,
				}, nil)
			},
			check: func() {
				// Check cache
				// `ca` is from the packageKey of cache library
				key := fmt.Sprintf("ca:%s:%s", domain.CACHE_PREFIX_SHORT_URL, "testid1")
				b, err := s.ring.Get(s.ctx, key).Bytes()
				s.NoError(err)

				var obj domain.GetRespDto
				s.NoError(json.Unmarshal(b, &obj))
				s.Equal(&domain.GetRespDto{
					Status:   domain.GetRespStatusNormal,
					Url:      "https://example.com/whatever1",
					ExpireAt: s.now,
				}, &obj)
			},
			expObj: &domain.GetRespDto{
				Status:   domain.GetRespStatusNormal,
				Url:      "https://example.com/whatever1",
				ExpireAt: s.now,
			},
			expErr: nil,
		},
		{
			name: "failed to get record when the record not found",
			req:  "testid1",
			setup: func() {
				s.repo.On("Get", s.ctx, "testid1").Once().Return(nil, domain.ErrRecordNotFound)
			},
			check: func() {
				// Check cache
				// `ca` is from the packageKey of cache library
				key := fmt.Sprintf("ca:%s:%s", domain.CACHE_PREFIX_SHORT_URL, "testid1")
				b, err := s.ring.Get(s.ctx, key).Bytes()
				s.NoError(err)

				var obj domain.GetRespDto
				s.NoError(json.Unmarshal(b, &obj))
				s.Equal(&domain.GetRespDto{Status: domain.GetRespStatusNotFound}, &obj)
			},
			expObj: &domain.GetRespDto{
				Status:   domain.GetRespStatusNotFound,
				Url:      "",
				ExpireAt: time.Time{},
			},
			expErr: nil,
		},
		{
			name: "failed to get record when the record is expired",
			req:  "testid1",
			setup: func() {
				s.repo.On("Get", s.ctx, "testid1").Once().Return(&domain.GetRespDto{
					Url:      "https://example.com/whatever1",
					ExpireAt: s.now.Add(-1 * time.Second),
				}, nil)
			},
			check: func() {
				// Check cache
				// `ca` is from the packageKey of cache library
				key := fmt.Sprintf("ca:%s:%s", domain.CACHE_PREFIX_SHORT_URL, "testid1")
				b, err := s.ring.Get(s.ctx, key).Bytes()
				s.NoError(err)

				var obj domain.GetRespDto
				s.NoError(json.Unmarshal(b, &obj))
				s.Equal(&domain.GetRespDto{
					Status:   domain.GetRespStatusNormal,
					Url:      "https://example.com/whatever1",
					ExpireAt: s.now.Add(-1 * time.Second),
				}, &obj)
			},
			expObj: &domain.GetRespDto{
				Status:   domain.GetRespStatusExpired,
				Url:      "https://example.com/whatever1",
				ExpireAt: s.now.Add(-1 * time.Second),
			},
			expErr: nil,
		},
		{
			name: "failed to get record due to unknown error",
			req:  "testid2",
			setup: func() {
				s.repo.On("Get", s.ctx, "testid2").Once().Return(nil, errors.New("unknown error"))
			},
			expObj: nil,
			expErr: errors.New("unknown error"),
		},
	} {
		s.Suite.Run(t.name, func() {
			ctx := context.Background()
			if t.setup != nil {
				t.setup()
			}
			obj, err := s.impl.Get(ctx, t.req)
			s.Equal(t.expErr, err)
			s.Equal(t.expObj, obj)
			if t.check != nil {
				t.check()
			}

		})
	}
}

func ConnectToDockerTestRedis() (string, string, func() error, error) {
	// Set up test db
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", "", nil, fmt.Errorf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return "", "", nil, fmt.Errorf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("redis", "7.4-alpine", []string{"TZ=UTC"})
	if err != nil {
		return "", "", nil, fmt.Errorf("Could not start resource: %s", err)
	}

	port := resource.GetPort("6379/tcp")
	host := "localhost"
	addr := fmt.Sprintf("%s:%s", host, port)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{Addr: addr})
		defer client.Close()
		_, e := client.Ping(context.Background()).Result()
		log.Printf("ping to redis(%s)", addr)
		return e
	}); err != nil {
		return "", "", nil, fmt.Errorf("Could not connect to database: %s", err)
	}

	return host, port, func() error {
		return pool.Purge(resource)
	}, nil
}
