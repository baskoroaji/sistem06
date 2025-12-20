package entity

import "errors"

var (
	ErrUserNotPermitted = errors.New("User Not Permitted")
	ErrUnauthorized     = errors.New("User Unathorized")
)
