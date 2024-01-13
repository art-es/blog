package v1_user_activate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
)

const invalidCodeErrorMessage = "invalid code"

type usecase interface {
	Do(ctx context.Context, in *dto.ActivateIn) error
}

type response struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type Handler struct {
	usecase            usecase
	serverErrorHandler api.ServerErrorHandler
}

func New(usecase usecase, serverErrorHandler api.ServerErrorHandler) *Handler {
	return &Handler{
		usecase:            usecase,
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
	if code == "" || !validUUID(code) {
		ctx.JSON(http.StatusBadRequest, response{Message: invalidCodeErrorMessage})
		return
	}

	err := h.usecase.Do(ctx, &dto.ActivateIn{Code: code})
	if err != nil {
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

func validUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
