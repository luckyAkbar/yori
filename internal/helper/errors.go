package helper

import "errors"

var (
	ErrInvalidHeaderErr      = errors.New("invalid header")
	ErrInvalidValueToRecords = errors.New("invalid value")
)
