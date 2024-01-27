package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type accessTokenParser interface {
	ParseAndValidate(token string) (*AccessTokenObject, error)
}

type AccessTokenParseCase struct {
	userRepository    UserRepository
	accessTokenParser accessTokenParser
}

func NewAccessTokenParseCase(
	userRepository UserRepository,
	accessTokenParser accessTokenParser,
) *AccessTokenParseCase {
	return &AccessTokenParseCase{
		userRepository:    userRepository,
		accessTokenParser: accessTokenParser,
	}
}

func (u *AccessTokenParseCase) Use(ctx context.Context, in *dto.AccessTokenParseIn) (*dto.ParseTokenOut, error) {
	tokenObject, err := u.accessTokenParser.ParseAndValidate(in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("access token validation error: %w", err)
	}

	if err = checkUserExistence(u.userRepository, ctx, tokenObject.UserID); err != nil {
		return nil, err
	}

	return &dto.ParseTokenOut{UserID: tokenObject.UserID}, nil
}
