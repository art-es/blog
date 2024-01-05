package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const messageInternalServerError = "sorry, try again later"

type errorResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type ErrorCatcher struct {
	logger      *zap.Logger
	handlerName string
}

func NewErrorCatcher(logger *zap.Logger, handler string) *ErrorCatcher {
	return &ErrorCatcher{
		logger:      logger,
		handlerName: handler,
	}
}

func (c *ErrorCatcher) Catch(ctx *gin.Context, err error) {
	c.logger.Error("API internal server error",
		zap.Error(err),
		zap.String("handler", c.handlerName))

	ctx.JSON(http.StatusInternalServerError, errorResponse{
		OK:      false,
		Message: messageInternalServerError,
	})
}
