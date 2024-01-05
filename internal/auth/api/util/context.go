package util

import "github.com/gin-gonic/gin"

const userIDCtxKey = "user_id"

func SetUserID(ctx *gin.Context, value int64) {
	ctx.Set(userIDCtxKey, value)
}

func GetUserID(ctx *gin.Context) int64 {
	return ctx.GetInt64(userIDCtxKey)
}
