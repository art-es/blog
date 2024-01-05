package dto

import "errors"

var (
	ErrEmailIsBusy            = errors.New("email is busy")
	ErrUserNotFound           = errors.New("auth not found")
	ErrActivationCodeNotFound = errors.New("activation code not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrInvalidAccessToken     = errors.New("invalid access token")
)
