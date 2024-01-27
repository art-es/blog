package v1_user_register

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
	Name     string `json:"name" validate:"required,lte=255"`
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=70"`
}

type response struct {
	Message string `json:"message,omitempty"`
}

type handler struct {
	userRegisterCase   userRegisterCase
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
		case dto.ErrEmailIsBusy:
			auth_api.BusyEmailResponse(ctx)
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
	in := dto.UserRegisterIn{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	return h.userRegisterCase.Use(ctx, &in)
}

func okResponse(ctx *gin.Context) {
	const message = "Please check your email to activate your account."
	ctx.JSON(http.StatusOK, &response{Message: message})
}
