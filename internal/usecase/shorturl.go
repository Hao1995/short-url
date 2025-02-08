package usecase

import (
	"context"
	"fmt"
	"hash/crc32"
	"log"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/caarlos0/env/v11"
)

var cfg config

type config struct {
	AppHost string `env:"APP_HOST" envDefault:"http://localhost"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse env: ", err)
	}
}

var (
	now = func() time.Time {
		return time.Now().UTC()
	}
)

type ShortUrlUseCase struct {
	repo Repository
}

// NewShortUrlUseCase generates the use case implementation of the ShortUrl use case interface
func NewShortUrlUseCase(repo Repository) UseCase {
	return &ShortUrlUseCase{
		repo: repo,
	}
}

// Create creates short_url record and return short url id
func (uc *ShortUrlUseCase) Create(ctx context.Context, CreateReqDto *domain.CreateReqDto) (*domain.CreateRespDto, error) {
	CreateReqDto.TargetID = fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(CreateReqDto.Url)))
	id, err := uc.repo.Create(ctx, CreateReqDto)
	if err != nil {
		return nil, err
	}
	return &domain.CreateRespDto{
		TargetID: id,
		ShortUrl: fmt.Sprintf("%s/%s", cfg.AppHost, id),
	}, nil
}

// Get gets short url record by id
func (uc *ShortUrlUseCase) Get(ctx context.Context, id string) (*domain.GetRespDto, error) {
	obj, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if obj.ExpiredAt.Before(now()) {
		return nil, domain.ErrExpired
	}

	return obj, nil
}
