package authenticated

import (
	"net/http"

	"github.com/art-es/blog/internal/common/apiutil"
	"github.com/gin-gonic/gin"
)

type response struct {
	OK bool `json:"ok"`
}

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Handle(ctx *gin.Context) {
	switch apiutil.GetUserID(ctx) {
	case 0:
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response{})
	default:
		ctx.Next()
	}
}
