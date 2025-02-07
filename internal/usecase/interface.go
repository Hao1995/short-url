package usecase

import (
	"context"

	"github.com/Hao1995/short-url/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, createDto *domain.CreateDto) (string, error)
	Get(ctx context.Context, id string) (string, error)
}

type UseCase interface {
	Create(ctx context.Context, createDto *domain.CreateDto) (string, error)
	Get(ctx context.Context, id string) (string, error)
}
