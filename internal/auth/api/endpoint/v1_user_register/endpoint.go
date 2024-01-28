//go:generate mockgen -source=endpoint.go -destination=mock/endpoint.go -package=mock
package v1_user_register

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/validation"
)

const (
	method = http.MethodPost
	path   = "/v1/auth/user/register"
)

type userRegisterCase interface {
	Use(ctx context.Context, in *dto.UserRegisterIn) error
}

func Bind(
	router *gin.Engine,
	userRegisterCase userRegisterCase,
	validator validation.Validator,
	serverErrorHandlerFactory api.ServerErrorHandlerFactory,
) {
	h := handler{
		userRegisterCase:   userRegisterCase,
		validator:          validator,
		serverErrorHandler: serverErrorHandlerFactory.MakeHandler(method, path),
	}

	router.Handle(method, path, h.handle)
}
