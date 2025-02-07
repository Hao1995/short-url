package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/internal/usecase"
	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_PASSWORD    = "password"
	MIGRATION_PATH = "../../../../cmd/shorturl/migration"
)

type ShortUrlTestSuite struct {
	suite.Suite
	dockertestClose func() error
	db              *gorm.DB

	impl usecase.Repository
}

func TestShortUrlTestSuite(t *testing.T) {
	suite.Run(t, new(ShortUrlTestSuite))
}

func (s *ShortUrlTestSuite) SetupSuite() {
	var err error
	var dbDSN string
	dbDSN, s.dockertestClose, err = ConnectToDockerTestDB()
	if err != nil {
		log.Fatal("failed to connect to docker test DB", err)
	}

	if err := GooseMigrate(dbDSN); err != nil {
		log.Fatal("failed to migrate DB", err)
	}

	s.db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("failed to init GORM connection", err)
	}

	s.impl = NewShortUrlRepository(s.db)

	now = func() time.Time {
		return time.Date(2025, 2, 10, 8, 30, 15, 0, time.UTC)
	}
}

func (s *ShortUrlTestSuite) SetupTest() {}

func (s *ShortUrlTestSuite) TearDownSubTest() {
	s.db.Where("1=1").Delete(&ShortUrl{})
}

func (s *ShortUrlTestSuite) TearDownTest() {}

func (s *ShortUrlTestSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	s.dockertestClose()
}

func (s *ShortUrlTestSuite) TestCreate() {
	for _, t := range []struct {
		name   string
		setup  func()
		req    *domain.CreateDto
		expID  string
		expErr error
	}{
		{
			name: "create record successfully",
			req: &domain.CreateDto{
				Url:       "https://example.com/whatever1",
				TargetID:  "testid1",
				ExpiredAt: now(),
			},
			expID:  "testid1",
			expErr: nil,
		},
		{
			name: "failed to create record due to duplicated target_id",
			setup: func() {
				shortUrl := ShortUrl{
					Url:       "https://example.com/whatever1",
					TargetID:  "testid1",
					ExpiredAt: now(),
					CreatedAt: now(),
				}
				s.Suite.Nil(s.db.Create(&shortUrl).Error)
			},
			req: &domain.CreateDto{
				Url:       "https://example.com/whatever2",
				TargetID:  "testid1",
				ExpiredAt: now(),
			},
			expID:  "",
			expErr: ErrDuplicatedKey,
		},
	} {
		s.Suite.Run(t.name, func() {
			ctx := context.Background()
			if t.setup != nil {
				t.setup()
			}
			id, err := s.impl.Create(ctx, t.req)
			s.ErrorIs(err, t.expErr)
			s.Equal(t.expID, id)
		})
	}
}

func (s *ShortUrlTestSuite) TestGet() {
	for _, t := range []struct {
		name   string
		req    string
		setup  func()
		expUrl string
		expErr error
	}{
		{
			name: "get record successfully",
			setup: func() {
				shortUrl := ShortUrl{
					Url:       "https://example.com/whatever1",
					TargetID:  "testid1",
					ExpiredAt: now(),
					CreatedAt: now(),
				}
				s.Suite.Nil(s.db.Create(&shortUrl).Error)
			},
			req:    "testid1",
			expUrl: "https://example.com/whatever1",
			expErr: nil,
		},
		{
			name:   "record not found",
			req:    "testid1",
			expUrl: "",
			expErr: ErrRecordNotFound,
		},
	} {
		s.Suite.Run(t.name, func() {
			ctx := context.Background()
			if t.setup != nil {
				t.setup()
			}
			url, err := s.impl.Get(ctx, t.req)
			s.ErrorIs(err, t.expErr)
			s.Equal(t.expUrl, url)
		})
	}
}

func ConnectToDockerTestDB() (string, func() error, error) {
	// Set up test db
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", nil, fmt.Errorf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return "", nil, fmt.Errorf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0", []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", DB_PASSWORD)})
	if err != nil {
		return "", nil, fmt.Errorf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var dsn string
	if err := pool.Retry(func() error {
		// https://gorm.io/docs/connecting_to_the_database.html#MySQL
		var err error
		dsn = fmt.Sprintf("root:%s@tcp(localhost:%s)/mysql?charset=utf8mb4&parseTime=True&loc=Local", DB_PASSWORD, resource.GetPort("3306/tcp"))
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		return "", nil, fmt.Errorf("Could not connect to database: %s", err)
	}

	return dsn, func() error {
		return pool.Purge(resource)
	}, nil
}

func GooseMigrate(dbString string) error {
	db, err := goose.OpenDBWithDriver("mysql", dbString)
	if err != nil {
		return fmt.Errorf("sql connection failed: %s", err)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, "up", db, MIGRATION_PATH); err != nil {
		return fmt.Errorf("goose up: %v", err)
	}

	return nil
}
