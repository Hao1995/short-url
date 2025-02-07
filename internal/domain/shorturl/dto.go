package shorturl

import "time"

type CreateDto struct {
	Url       string
	ExpiredAt time.Time
}
