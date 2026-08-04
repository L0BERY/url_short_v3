package repository

import "errors"

var (
	ErrCodeNotFound = errors.New("repository: code not found")
	ErrUrlNotFound  = errors.New("repository: url not found")
)
