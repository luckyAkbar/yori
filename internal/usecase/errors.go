package usecase

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request error")
	ErrInternal   = errors.New("internal error")
	ErrInProgress = errors.New("in progress operation")
)
