package domain

import "errors"

var (
	ErrDuplicatedKey  = errors.New("duplicated key")
	ErrRecordNotFound = errors.New("record not found")
)
