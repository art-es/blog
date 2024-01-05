package domain

import "time"

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	Active       bool
}

type AccessTokenObject struct {
	ExpirationTime time.Time
	NotBefore      time.Time
	IssuedAt       time.Time
	Audience       []string
	Issuer         string
	UserID         int64
}
