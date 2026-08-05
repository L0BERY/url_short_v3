package service

import "errors"

var (
	ErrUrlNotFound = errors.New("service: url not found")
	ErrEmptyCode   = errors.New("service: code is empty")
	ErrEmptyUrl    = errors.New("service: url is empty")

	ErrTooManyAttempts = errors.New("service: err too many attempts")

	ErrServer = errors.New("service: internal error")
)
