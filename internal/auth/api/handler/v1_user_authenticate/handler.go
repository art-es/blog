package v1_user_authenticate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
)

const credentialsErrorMessage = "credentials are wrong"

type usecase interface {
	Do(ctx context.Context, in *dto.AuthenticateIn) (*dto.AuthenticateOut, error)
}

type request struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=70"`
}

type response struct {
	OK          bool   `json:"ok"`
	AccessToken string `json:"accessToken,omitempty"`
	Message     string `json:"message,omitempty"`
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
	return "/v1/auth/user/authenticate"
}

func (h *Handler) Handle(ctx *gin.Context) {
	var req request
	_ = ctx.BindJSON(&req)

	if err := h.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
		return
	}

	in := dto.AuthenticateIn(req)
	out, err := h.usecase.Do(ctx, &in)
	if err != nil {
		switch err {
		case dto.ErrUserNotFound, dto.ErrWrongPassword:
			ctx.JSON(http.StatusBadRequest, response{Message: credentialsErrorMessage})
		default:
			h.serverErrorHandler.Handle(h.Endpoint(), ctx.Writer, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, response{
		OK:          true,
		AccessToken: out.AccessToken,
	})
}
