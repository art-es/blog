package apiutil

import "github.com/gin-gonic/gin"

const userIDContextKey = "user_id"

func SetUserID(ctx *gin.Context, value int64) {
	ctx.Set(userIDContextKey, value)
}

func GetUserID(ctx *gin.Context) int64 {
	return ctx.GetInt64(userIDContextKey)
}
