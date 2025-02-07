package usecase

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/Hao1995/short-url/internal/domain"
)

type ShortUrlUseCase struct {
	repo Repository
}

func NewShortUrlUseCase(repo Repository) UseCase {
	return &ShortUrlUseCase{
		repo: repo,
	}
}

func (uc *ShortUrlUseCase) Create(ctx context.Context, createDto *domain.CreateDto) (string, error) {
	createDto.TargetID = fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(createDto.Url)))
	id, err := uc.repo.Create(ctx, createDto)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (uc *ShortUrlUseCase) Get(ctx context.Context, id string) (string, error) {
	id, err := uc.repo.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return id, nil
}
