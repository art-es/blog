package v1_user_activate

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
	Code string `json:"code" validate:"required,uuid"`
}

type response struct {
	Message string `json:"message,omitempty"`
}

type handler struct {
	userActivateCase   userActivateCase
	validator          validation.Validator
	serverErrorHandler api.ServerErrorHandler
}

func (h *handler) handle(ctx *gin.Context) {
	req, err := h.parseRequest(ctx)
	if err != nil {
		api.RequestValidationFailedResponse(ctx, err)
		return
	}

	if err = h.useCase(ctx, req); err != nil {
		switch err {
		case dto.ErrUserActivationCodeNotFound, dto.ErrUserNotFound:
			notFoundResponse(ctx)
		case dto.ErrExpiredUserActivationCode:
			auth_api.ExpiredUserActivationCodeResponse(ctx)
		default:
			h.serverErrorHandler.Handle(ctx, err)
		}
		return
	}

	okResponse(ctx)
}

func (h *handler) parseRequest(ctx *gin.Context) (*request, error) {
	var req request
	ctx.ShouldBindJSON(&req)

	if err := h.validator.Struct(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (h *handler) useCase(ctx context.Context, req *request) error {
	in := dto.UserActivateIn{
		Code: req.Code,
	}

	return h.userActivateCase.Use(ctx, &in)
}

func notFoundResponse(ctx *gin.Context) {
	const message = "Activation code not found."
	ctx.JSON(http.StatusNotFound, &response{Message: message})
}

func okResponse(ctx *gin.Context) {
	const message = "Your account have been activated."
	ctx.JSON(http.StatusOK, &response{Message: message})
}
