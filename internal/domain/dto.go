package domain

import "time"

type CreateReqDto struct {
	Url       string
	TargetID  string
	ExpiredAt time.Time
}

type GetRespDto struct {
	Url       string
	ExpiredAt time.Time
}
