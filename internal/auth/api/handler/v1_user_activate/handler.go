//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package v1_user_activate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
)

type usecase interface {
	Do(ctx context.Context, in *dto.ActivateIn) error
}

type response struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type Handler struct {
	usecase            usecase
	validator          api.Validator
	serverErrorHandler api.ServerErrorHandler
}

func New(
	usecase usecase,
	validator api.Validator,
	serverErrorHandler api.ServerErrorHandler,
) *Handler {
	return &Handler{
		usecase:            usecase,
		validator:          validator,
		serverErrorHandler: serverErrorHandler,
	}
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Endpoint() string {
	return "/v1/auth/user/activate/:code"
}

func (h *Handler) Handle(ctx *gin.Context) {
	code := ctx.Param("code")

	if err := h.validator.Var(code, "required,uuid"); err != nil {
		ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
		return
	}

	in := dto.ActivateIn{Code: code}
	if err := h.usecase.Do(ctx, &in); err != nil {
		switch err {
		case dto.ErrActivationCodeNotFound, dto.ErrUserNotFound:
			ctx.JSON(http.StatusNotFound, response{})
		default:
			h.serverErrorHandler.Handle(h.Endpoint(), ctx.Writer, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response{OK: true})
}
