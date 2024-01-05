package register

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/art-es/blog/internal/auth/dto"
)

const internalServerErrorMessage = "sorry, try again later"

var validate = validator.New()

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
	usecase usecase
	logger  *zap.Logger
}

func NewHandler(usecase usecase, logger *zap.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Handle(ctx *gin.Context) {
	var req request
	_ = ctx.BindJSON(&req)

	if err := validate.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response{OK: false, Message: err.Error()})
		return
	}

	in := dto.RegisterIn(req)
	if err := h.usecase.Do(ctx, &in); err != nil {
		switch err {
		case dto.ErrEmailIsBusy:
			ctx.JSON(http.StatusBadRequest, response{OK: false, Message: err.Error()})
		default:
			h.logger.Error("auth: register api error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, response{OK: false, Message: internalServerErrorMessage})
		}
		return
	}

	ctx.JSON(http.StatusOK, response{OK: true})
}
