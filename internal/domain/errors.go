package domain

import "errors"

var (
	ErrDuplicatedKey  = errors.New("duplicated key")
	ErrExpired        = errors.New("expired")
	ErrRecordNotFound = errors.New("record not found")
)
