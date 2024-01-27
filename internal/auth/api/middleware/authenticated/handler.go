package authenticated

import (
	"github.com/art-es/blog/internal/common/api"

	"github.com/gin-gonic/gin"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handle(ctx *gin.Context) {
	switch api.GetUserID(ctx) {
	case 0:
		api.UnauthorizedResponse(ctx)
	default:
		ctx.Next()
	}
}
