package access_token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Registered Claim Names of JSON Web AccessToken
// https://datatracker.ietf.org/doc/html/rfc7519#section-4
type claims struct {
	ExpirationTime *jwt.NumericDate  `json:"exp"`
	NotBefore      *jwt.NumericDate  `json:"nbf"`
	IssuedAt       *jwt.NumericDate  `json:"iat"`
	Audience       *jwt.ClaimStrings `json:"aud"`
	Issuer         string            `json:"iss"`
	Subject        string            `json:"sub"`
}

func (c *claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.ExpirationTime, nil
}

func (c *claims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.NotBefore, nil
}

func (c *claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.IssuedAt, nil
}

func (c *claims) GetAudience() (jwt.ClaimStrings, error) {
	if c.Audience == nil {
		return nil, errors.New("aud is nil")
	}
	return *c.Audience, nil
}

func (c *claims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c *claims) GetSubject() (string, error) {
	return c.Subject, nil
}
