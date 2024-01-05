package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type accessTokenParser interface {
	ParseAndValidate(token string) (*AccessTokenObject, error)
}

type ParseTokenUsecase struct {
	userRepository    UserRepository
	accessTokenParser accessTokenParser
}

func NewParseTokenUsecase(
	userRepository UserRepository,
	accessTokenParser accessTokenParser,
) *ParseTokenUsecase {
	return &ParseTokenUsecase{
		userRepository:    userRepository,
		accessTokenParser: accessTokenParser,
	}
}

func (u *ParseTokenUsecase) Do(ctx context.Context, in *dto.ParseTokenIn) (*dto.ParseTokenOut, error) {
	tokenObject, err := u.accessTokenParser.ParseAndValidate(in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("access token validation error: %w", err)
	}

	if err = checkUserExistence(u.userRepository, ctx, tokenObject.UserID); err != nil {
		return nil, err
	}

	return &dto.ParseTokenOut{UserID: tokenObject.UserID}, nil
}
