package request

import "time"

type ShortUrlCreateRequest struct {
	Url       string    `form:"url" json:"url" binding:"required"`
	ExpiredAt time.Time `form:"expiredAt" json:"expiredAt" binding:"required"`
}

type ShortUrlGetRequest struct {
	ID string `uri:"id" binding:"required"`
}
