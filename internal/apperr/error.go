package apperr

import "errors"

var (
	ErrUserIsBlocked = errors.New("user is blocked")
)
