package gotv

import "errors"

var (
	ErrInvalidAuth      = errors.New("Invalid Authentication")
	ErrFragmentNotFound = errors.New("Fragment Not Found")
	ErrMatchNotFound    = errors.New("Match Not Found")
)
