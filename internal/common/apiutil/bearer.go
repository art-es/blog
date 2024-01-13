package apiutil

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseBearerToken(ctx *gin.Context) (string, bool) {
	values := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(values) != 2 || strings.ToLower(values[0]) != "bearer" || values[1] == "" {
		return "", false
	}
	return values[1], true
}
