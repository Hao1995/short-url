package usecase

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/pkg/migrationkit/randkit"
	"github.com/caarlos0/env/v11"
	"github.com/viney-shih/go-cache"
)

var (
	now = func() time.Time {
		return time.Now()
	}

	randString = func() string {
		return randkit.String(4)
	}

	cfg config
)

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse env: ", err)
	}
}

type config struct {
	AppHost string `env:"APP_HOST" envDefault:"http://localhost"`
}
type ShortUrlUseCase struct {
	repo Repository
	c    cache.Cache
}

// NewShortUrlUseCase generates the use case implementation of the ShortUrl use case interface
func NewShortUrlUseCase(repo Repository, c cache.Cache) UseCase {
	return &ShortUrlUseCase{
		repo: repo,
		c:    c,
	}
}

// Create creates short_url record and return short url id
func (uc *ShortUrlUseCase) Create(ctx context.Context, createReqDto *domain.CreateReqDto) (*domain.CreateRespDto, error) {
	var id string
	url := createReqDto.Url
	for {
		var err error
		createReqDto.TargetID = fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(url)))
		id, err = uc.repo.Create(ctx, createReqDto)
		if err == domain.ErrDuplicatedKey {
			url += randString() // 62^4=14M possibilities
			log.Print("Append random suffix", url)
		} else if err != nil {
			return nil, err
		} else {
			break
		}
	}
	return &domain.CreateRespDto{
		TargetID: id,
		ShortUrl: fmt.Sprintf("%s/%s", cfg.AppHost, id),
	}, nil
}

// Get gets short url record by id
func (uc *ShortUrlUseCase) Get(ctx context.Context, id string) (*domain.GetRespDto, error) {
	cacheObj := &domain.GetRespDto{}
	if err := uc.c.GetByFunc(ctx, domain.CACHE_PREFIX_SHORT_URL, id, cacheObj, func() (interface{}, error) {
		obj, err := uc.repo.Get(ctx, id)
		if err == domain.ErrRecordNotFound {
			obj = &domain.GetRespDto{Status: domain.GetRespStatusNotFound}
		} else if err != nil {
			return nil, err
		} else {
			obj.Status = domain.GetRespStatusNormal
		}
		return obj, nil
	}); err != nil {
		log.Print("ShortUrlUseCase.Get. Failed to get the short_url from cache: ", err)
		return nil, err
	}

	if cacheObj.Status == domain.GetRespStatusNormal {
		if cacheObj.ExpireAt.Before(now()) {
			cacheObj.Status = domain.GetRespStatusExpired
		}
	}

	return cacheObj, nil
}
