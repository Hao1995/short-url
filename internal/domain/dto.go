//go:generate go-enum --marshal
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

// ENUM(Normal, NotFound, Expired)
type GetRespStatus string

type GetRespDto struct {
	Status    GetRespStatus
	Url       string
	ExpiredAt time.Time
}
