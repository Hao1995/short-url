package usecase

import (
	"context"

	"github.com/Hao1995/short-url/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, CreateReqDto *domain.CreateReqDto) (string, error)
	Get(ctx context.Context, id string) (*domain.GetRespDto, error)
}

type UseCase interface {
	Create(ctx context.Context, CreateReqDto *domain.CreateReqDto) (*domain.CreateRespDto, error)
	Get(ctx context.Context, id string) (*domain.GetRespDto, error)
}
