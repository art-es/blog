//go:generate mockgen -source=case_refresh_access_token.go -destination=mock/case_refresh_access_token.go -package=mock
package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type accessTokenRefresher interface {
	Parse(token string) (*AccessTokenObject, error)
	Refresh(object *AccessTokenObject)
	Sign(object *AccessTokenObject) (string, error)
}

type AccessTokenRefreshCase struct {
	userRepository       UserRepository
	accessTokenRefresher accessTokenRefresher
}

func NewAccessTokenRefreshCase(
	userRepository UserRepository,
	accessTokenRefresher accessTokenRefresher,
) *AccessTokenRefreshCase {
	return &AccessTokenRefreshCase{
		userRepository:       userRepository,
		accessTokenRefresher: accessTokenRefresher,
	}
}

func (c *AccessTokenRefreshCase) Use(ctx context.Context, in *dto.AccessTokenRefreshIn) (*dto.AccessTokenRefreshOut, error) {
	tokenObject, err := c.accessTokenRefresher.Parse(in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("access token parsing error: %w", err)
	}

	if err = checkUserExistence(c.userRepository, ctx, tokenObject.UserID); err != nil {
		return nil, err
	}

	token, err := refreshAccessToken(c.accessTokenRefresher, tokenObject)
	if err != nil {
		return nil, err
	}

	return &dto.AccessTokenRefreshOut{AccessToken: token}, nil
}

func refreshAccessToken(refresher accessTokenRefresher, object *AccessTokenObject) (string, error) {
	refresher.Refresh(object)

	token, err := refresher.Sign(object)
	if err != nil {
		return "", fmt.Errorf("signing access token object error: %w", err)
	}

	return token, nil
}

func checkUserExistence(userRepository UserRepository, ctx context.Context, userID int64) error {
	exists, err := userRepository.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("user checking existence in repository error: %w", err)
	}
	if !exists {
		return dto.ErrUserNotFound
	}
	return nil
}
