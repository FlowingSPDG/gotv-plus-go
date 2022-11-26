package gotv

import "golang.org/x/xerrors"

var (
	ErrInvalidAuth      = xerrors.New("Invalid Authentication")
	ErrFragmentNotFound = xerrors.New("Fragment Not Found")
	ErrMatchNotFound    = xerrors.New("Match Not Found")
)
