//go:generate mockgen -source=server_error.go -destination=mock/server_error.go -package=mock
package api

import (
	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/common/log"
)

type ServerErrorHandlerFactory interface {
	MakeHandler(method, path string) ServerErrorHandler
}

type ServerErrorHandler interface {
	Handle(ctx *gin.Context, err error)
}

type serverErrorHandlerFactory struct {
	logger log.Logger
}

type serverErrorHandler struct {
	endpoint string
	logger   log.Logger
}

func NewServerErrorHandlerFactory(logger log.Logger) ServerErrorHandlerFactory {
	return &serverErrorHandlerFactory{
		logger: logger,
	}
}

func (f *serverErrorHandlerFactory) MakeHandler(method, path string) ServerErrorHandler {
	return &serverErrorHandler{
		endpoint: method + " " + path,
		logger:   f.logger,
	}
}

func (h *serverErrorHandler) Handle(ctx *gin.Context, err error) {
	h.logger.Error("API server error",
		log.String("endpoint", h.endpoint),
		log.Error(err),
	)

	InternalServerErrorResponse(ctx)
}
