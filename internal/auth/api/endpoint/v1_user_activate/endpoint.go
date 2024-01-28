//go:generate mockgen -source=endpoint.go -destination=mock/endpoint.go -package=mock
package v1_user_activate

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
	path   = "/v1/auth/user/activate"
)

type userActivateCase interface {
	Use(ctx context.Context, in *dto.UserActivateIn) error
}

func Bind(
	router *gin.Engine,
	userActivateCase userActivateCase,
	validator validation.Validator,
	serverErrorHandlerFactory api.ServerErrorHandlerFactory,
) {
	h := handler{
		userActivateCase:   userActivateCase,
		validator:          validator,
		serverErrorHandler: serverErrorHandlerFactory.MakeHandler(method, path),
	}

	router.Handle(method, path, h.handle)
}
