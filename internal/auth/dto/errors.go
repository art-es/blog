package dto

import "errors"

var (
	ErrEmailIsBusy                = errors.New("busy email")
	ErrUserNotFound               = errors.New("auth not found")
	ErrUserActivationCodeNotFound = errors.New("activation code not found")
	ErrExpiredUserActivationCode  = errors.New("expired user activation code")
	ErrIncorrectPassword          = errors.New("incorrect password")
	ErrInvalidAccessToken         = errors.New("invalid access token")
)
