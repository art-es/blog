//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package v1_user_register

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
)

type usecase interface {
	Do(ctx context.Context, in *dto.RegisterIn) error
}

type request struct {
	Name     string `json:"name" validate:"required,lte=255"`
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=70"`
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
	return http.MethodPost
}

func (h *Handler) Endpoint() string {
	return "/v1/auth/user/register"
}

func (h *Handler) Handle(ctx *gin.Context) {
	var req request
	ctx.ShouldBindJSON(&req)

	if err := h.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
		return
	}

	in := dto.RegisterIn(req)
	if err := h.usecase.Do(ctx, &in); err != nil {
		switch err {
		case dto.ErrEmailIsBusy:
			ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
		default:
			h.serverErrorHandler.Handle(h.Endpoint(), ctx.Writer, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response{OK: true})
}
