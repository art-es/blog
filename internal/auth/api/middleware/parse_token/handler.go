//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package parse_token

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/apiutil"
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
	if accessToken, ok := apiutil.ParseBearerToken(ctx); ok {
		if out, err := m.usecase.Do(ctx, &dto.ParseTokenIn{AccessToken: accessToken}); err == nil {
			apiutil.SetUserID(ctx, out.UserID)
		}
	}

	ctx.Next()
}
