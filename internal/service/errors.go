package service

import "errors"

var (
	ErrUrlNotFound = errors.New("service: url not found")
	ErrEmptyCode   = errors.New("service: code if empty")

	ErrServer = errors.New("service: internal error")
)
