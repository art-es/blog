package api

import "github.com/gin-gonic/gin"

func AccessTokenHeader(ctx *gin.Context) string {
	return ctx.GetHeader("X-Access-Token")
}
