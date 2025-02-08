package domain

import "time"

type CreateReqDto struct {
	Url       string
	TargetID  string
	ExpiredAt time.Time
}

type CreateRespDto struct {
	TargetID string
	ShortUrl string
}

type GetRespDto struct {
	Url       string
	ExpiredAt time.Time
}
