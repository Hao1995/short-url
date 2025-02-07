package domain

import "time"

type CreateDto struct {
	Url       string
	TargetID  string
	ExpiredAt time.Time
}
