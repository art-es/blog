//go:generate mockgen -source=endpoint.go -destination=mock/endpoint.go -package=mock
package v1_access_token_refresh

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
)

const (
	method = http.MethodPost
	path   = "/v1/auth/access-token/refresh"
)

type accessTokenRefreshCase interface {
	Use(ctx context.Context, in *dto.AccessTokenRefreshIn) (*dto.AccessTokenRefreshOut, error)
}

func Bind(
	router *gin.Engine,
	accessTokenRefreshCase accessTokenRefreshCase,
	serverErrorHandlerFactory api.ServerErrorHandlerFactory,
) {
	h := handler{
		accessTokenRefreshCase: accessTokenRefreshCase,
		serverErrorHandler:     serverErrorHandlerFactory.MakeHandler(method, path),
	}

	router.Handle(method, path, h.handle)
}
