package mysql

import (
	"time"
)

// ShortUrl represents as table `short_urls`.
type ShortUrl struct {
	ID        uint `gorm:"primaryKey, autoIncrement"`
	Url       string
	TargetID  string `gorm:"uniqueIndex"`
	ExpiredAt time.Time
	CreatedAt time.Time
}
