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
		return time.Now()
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

// Create creates short_url record
func (repo *ShortUrlRepository) Create(ctx context.Context, createDto *domain.CreateDto) (string, error) {
	record := ShortUrl{
		Url:       createDto.Url,
		TargetID:  createDto.TargetID,
		ExpiredAt: createDto.ExpiredAt,
		CreatedAt: now(),
	}

	if result := repo.db.Create(&record); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return "", domain.ErrDuplicatedKey
		}
		log.Printf("failed to create short_url: %s", result.Error)
		return "", result.Error
	}

	return createDto.TargetID, nil
}

func (repo *ShortUrlRepository) Get(ctx context.Context, id string) (string, error) {
	var record ShortUrl
	if result := repo.db.Where("target_id = ?", id).Select("url").First(&record); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", domain.ErrRecordNotFound
		}
		log.Printf("failed to get short_url by id(%s): %s", id, result.Error)
		return "", result.Error
	}
	log.Printf("get url `%s` by id `%s`", record.Url, id)
	return record.Url, nil
}
