package parse_token

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/api/util"
	"github.com/art-es/blog/internal/auth/dto"
)

type usecase interface {
	Do(ctx context.Context, in *dto.ParseTokenIn) (*dto.ParseTokenOut, error)
}

type Middleware struct {
	usecase usecase
}

func New(usecase usecase) *Middleware {
	return &Middleware{usecase: usecase}
}

func (m *Middleware) Handle(ctx *gin.Context) {
	if accessToken, ok := util.ParseBearerToken(ctx); ok {
		if out, err := m.usecase.Do(ctx, &dto.ParseTokenIn{AccessToken: accessToken}); err == nil {
			util.SetUserID(ctx, out.UserID)
		}
	}

	ctx.Next()
}
