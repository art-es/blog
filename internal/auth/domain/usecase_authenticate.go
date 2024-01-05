//go:generate mockgen -source=usecase_authenticate.go -destination=mock/usecase_authenticate.go -package=mock
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

type AuthenticateUsecase struct {
	userRepository    UserRepository
	passwordValidator passwordValidator
	accessTokenIssuer accessTokenIssuer
}

func NewAuthenticateUsecase(
	userRepository UserRepository,
	passwordHashValidator passwordValidator,
	accessTokenService accessTokenIssuer,
) *AuthenticateUsecase {
	return &AuthenticateUsecase{
		userRepository:    userRepository,
		passwordValidator: passwordHashValidator,
		accessTokenIssuer: accessTokenService,
	}
}

func (u *AuthenticateUsecase) Do(ctx context.Context, in *dto.AuthenticateIn) (*dto.AuthenticateOut, error) {
	user, err := u.userRepository.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("auth getting by email error: %w", err)
	}
	if user == nil {
		return nil, dto.ErrUserNotFound
	}

	if err = u.passwordValidator.Validate(in.Password, user.PasswordHash); err != nil {
		return nil, fmt.Errorf("password validate by hash error: %w", err)
	}

	accessToken, err := newAccessToken(u.accessTokenIssuer, user.ID)
	if err != nil {
		return nil, fmt.Errorf("access token creation error: %w", err)
	}

	return &dto.AuthenticateOut{AccessToken: accessToken}, nil
}

func newAccessToken(issuer accessTokenIssuer, userID int64) (string, error) {
	object := issuer.NewObject(userID)
	return issuer.Sign(object)
}
