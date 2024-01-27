package v1_user_authenticate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	auth_api "github.com/art-es/blog/internal/auth/api"
	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/validation"
)

type request struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=70"`
}

type response struct {
	AccessToken string `json:"accessToken,omitempty"`
}

type handler struct {
	userAuthenticateCase userAuthenticateCase
	validator            validation.Validator
	serverErrorHandler   api.ServerErrorHandler
}

func (h *handler) handle(ctx *gin.Context) {
	req, err := h.parseRequest(ctx)
	if err != nil {
		api.RequestValidationFailedResponse(ctx, err)
		return
	}

	out, err := h.useCase(ctx, req)
	if err != nil {
		switch err {
		case dto.ErrUserNotFound, dto.ErrIncorrectPassword:
			auth_api.IncorrectUserCredentialsResponse(ctx)
		default:
			h.serverErrorHandler.Handle(ctx, err)
		}
		return
	}

	okResponse(ctx, out)
}

func (h *handler) parseRequest(ctx *gin.Context) (*request, error) {
	var req request
	ctx.ShouldBindJSON(&req)

	if err := h.validator.Struct(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *handler) useCase(ctx context.Context, req *request) (*dto.UserAuthenticateOut, error) {
	in := dto.UserAuthenticateIn{
		Email:    req.Email,
		Password: req.Password,
	}

	return h.userAuthenticateCase.Use(ctx, &in)
}

func okResponse(ctx *gin.Context, result *dto.UserAuthenticateOut) {
	ctx.JSON(http.StatusOK, &response{AccessToken: result.AccessToken})
}
