//go:generate mockgen -source=case_user_authenticate.go -destination=mock/case_user_authenticate.go -package=mock
package domain

import (
	"context"
	"fmt"

	"github.com/art-es/blog/internal/auth/dto"
)

type passwordValidator interface {
	Validate(password, hash string) error
}

type accessTokenIssuer interface {
	NewObject(userID int64) *AccessTokenObject
	Sign(object *AccessTokenObject) (string, error)
}

type UserAuthenticateCase struct {
	userRepository    UserRepository
	passwordValidator passwordValidator
	accessTokenIssuer accessTokenIssuer
}

func NewUserAuthenticateCase(
	userRepository UserRepository,
	passwordHashValidator passwordValidator,
	accessTokenService accessTokenIssuer,
) *UserAuthenticateCase {
	return &UserAuthenticateCase{
		userRepository:    userRepository,
		passwordValidator: passwordHashValidator,
		accessTokenIssuer: accessTokenService,
	}
}

func (c *UserAuthenticateCase) Use(ctx context.Context, in *dto.UserAuthenticateIn) (*dto.UserAuthenticateOut, error) {
	user, err := c.userRepository.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("auth getting by email error: %w", err)
	}
	if user == nil {
		return nil, dto.ErrUserNotFound
	}

	if err = c.passwordValidator.Validate(in.Password, user.PasswordHash); err != nil {
		return nil, fmt.Errorf("password validate by hash error: %w", err)
	}

	accessToken, err := newAccessToken(c.accessTokenIssuer, user.ID)
	if err != nil {
		return nil, fmt.Errorf("access token creation error: %w", err)
	}

	return &dto.UserAuthenticateOut{AccessToken: accessToken}, nil
}

func newAccessToken(issuer accessTokenIssuer, userID int64) (string, error) {
	object := issuer.NewObject(userID)
	return issuer.Sign(object)
}
