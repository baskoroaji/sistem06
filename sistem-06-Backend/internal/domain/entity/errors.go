package entity

import "errors"

var (
	ErrUserNotPermitted         = errors.New("User Not Permitted")
	ErrUnauthorized             = errors.New("User Unathorized")
	ErrDuplicateUser            = errors.New("user has already exist")
	ErrDuplicateUsernameOrEmail = errors.New("username and email has already exist")
)
