package models

import "errors"

var (
	ErrUserNotFound        = errors.New("user with these parameters not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
