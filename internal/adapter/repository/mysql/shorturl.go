package mysql

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Hao1995/short-url/internal/domain"
	"github.com/Hao1995/short-url/internal/usecase"
	"gorm.io/gorm"
)

var (
	now = func() time.Time {
		return time.Now().UTC()
	}
)

type ShortUrlRepository struct {
	db *gorm.DB
}

// NewShortUrlRepository generates the MySQL implementation of the ShortUrl repository interface
func NewShortUrlRepository(db *gorm.DB) usecase.Repository {
	return &ShortUrlRepository{
		db: db,
	}
}

// Create creates short_url record and return short url id
func (repo *ShortUrlRepository) Create(ctx context.Context, CreateReqDto *domain.CreateReqDto) (string, error) {
	record := ShortUrl{
		Url:       CreateReqDto.Url,
		TargetID:  CreateReqDto.TargetID,
		ExpiredAt: CreateReqDto.ExpiredAt,
		CreatedAt: now(),
	}

	if result := repo.db.Create(&record); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return "", domain.ErrDuplicatedKey
		}
		log.Printf("failed to create short_url: %s", result.Error)
		return "", result.Error
	}

	return CreateReqDto.TargetID, nil
}

// Get gets short url record by id
func (repo *ShortUrlRepository) Get(ctx context.Context, id string) (*domain.GetRespDto, error) {
	var record ShortUrl
	result := repo.db.Where("target_id = ?", id).Select([]string{"url", "expired_at"}).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		log.Printf("failed to get short_url by id(%s): %s", id, result.Error)
		return nil, result.Error
	}
	log.Printf("get url `%s` by id `%s`", record.Url, id)

	return &domain.GetRespDto{
		Url:       record.Url,
		ExpiredAt: record.ExpiredAt,
	}, nil
}
