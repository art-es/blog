//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package parse_token

import (
	"context"

	"github.com/art-es/blog/internal/common/api"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
)

type accessTokenParseCase interface {
	Use(ctx context.Context, in *dto.AccessTokenParseIn) (*dto.ParseTokenOut, error)
}

type Middleware struct {
	accessTokenParseCase accessTokenParseCase
}

func New(accessTokenParseCase accessTokenParseCase) *Middleware {
	return &Middleware{
		accessTokenParseCase: accessTokenParseCase,
	}
}

func (m *Middleware) Handle(ctx *gin.Context) {
	if accessToken := api.AccessTokenHeader(ctx); accessToken != "" {
		if userID, ok := m.useCase(ctx, accessToken); ok {
			api.SetUserID(ctx, userID)
		}
	}

	ctx.Next()
}

func (m *Middleware) useCase(ctx context.Context, accessToken string) (int64, bool) {
	in := dto.AccessTokenParseIn{
		AccessToken: accessToken,
	}

	if out, err := m.accessTokenParseCase.Use(ctx, &in); err == nil {
		return out.UserID, true
	}

	return 0, false
}
