package request

import "time"

type ShortUrlCreateRequest struct {
	Url      string    `form:"url" json:"url" binding:"required,url"`
	ExpireAt time.Time `form:"expireAt" json:"expireAt" binding:"required"`
}

type ShortUrlGetRequest struct {
	ID string `uri:"id" binding:"required"`
}
