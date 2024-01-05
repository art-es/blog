package access_token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/art-es/blog/internal/auth/domain"
)

type Service struct {
	secret              []byte
	expirationTimeShift time.Duration
	audience            []string
	issuer              string
	parserOpts          []jwt.ParserOption
}

func New(secret []byte) *Service {
	const (
		aud = "auth"
		iss = "art-es"
	)

	return &Service{
		secret:              secret,
		expirationTimeShift: 12 * time.Hour,
		audience:            []string{aud},
		issuer:              iss,
		parserOpts: []jwt.ParserOption{
			jwt.WithExpirationRequired(),
			jwt.WithIssuedAt(),
			jwt.WithAudience(aud),
			jwt.WithIssuer(iss),
		},
	}
}

func (s *Service) NewObject(userID int64) *domain.AccessTokenObject {
	now := time.Now()

	return &domain.AccessTokenObject{
		ExpirationTime: now.Add(s.expirationTimeShift),
		NotBefore:      now,
		IssuedAt:       now,
		Audience:       s.audience,
		Issuer:         s.issuer,
		UserID:         userID,
	}
}

func (s *Service) Refresh(object *domain.AccessTokenObject) {
	now := time.Now()

	object.ExpirationTime = now.Add(s.expirationTimeShift)
	object.NotBefore = now
	object.IssuedAt = now
}

func (s *Service) Sign(object *domain.AccessTokenObject) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, toClaims(object)).
		SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signing error: %w", err)
	}
	return token, nil
}

func (s *Service) Parse(token string) (*domain.AccessTokenObject, error) {
	clm := &claims{}

	_, err := jwt.ParseWithClaims(token, clm, s.parserFunc, s.parserOpts...)
	if err != nil {
		return nil, fmt.Errorf("parse string with secret error: %w", err)
	}

	obj, err := toAccessTokenObject(clm)
	if err != nil {
		return nil, fmt.Errorf("failed to convert claims to jwt object")
	}
	return obj, nil
}

func (s *Service) ParseAndValidate(token string) (*domain.AccessTokenObject, error) {
	clm := &claims{}
	parserOpts := append(s.parserOpts, jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithLeeway(2*time.Hour))

	_, err := jwt.ParseWithClaims(token, clm, s.parserFunc, parserOpts...)
	if err != nil {
		return nil, fmt.Errorf("parse string with secret error: %w", err)
	}

	obj, err := toAccessTokenObject(clm)
	if err != nil {
		return nil, fmt.Errorf("failed to convert claims to jwt object")
	}
	return obj, nil
}

func (s *Service) parserFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return s.secret, nil
}

func toClaims(o *domain.AccessTokenObject) *claims {
	aud := jwt.ClaimStrings(o.Audience)

	return &claims{
		ExpirationTime: jwt.NewNumericDate(o.ExpirationTime),
		NotBefore:      jwt.NewNumericDate(o.NotBefore),
		IssuedAt:       jwt.NewNumericDate(o.IssuedAt),
		Audience:       &aud,
		Issuer:         o.Issuer,
		Subject:        strconv.FormatInt(o.UserID, 10),
	}
}

func toAccessTokenObject(c *claims) (*domain.AccessTokenObject, error) {
	aud, err := c.GetAudience()
	if err != nil {
		return nil, err
	}

	userID, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse sub as int: %w", err)
	}

	return &domain.AccessTokenObject{
		ExpirationTime: c.ExpirationTime.Time,
		NotBefore:      c.NotBefore.Time,
		IssuedAt:       c.IssuedAt.Time,
		Audience:       aud,
		Issuer:         c.Issuer,
		UserID:         userID,
	}, nil
}
