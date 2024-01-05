package authenticate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/art-es/blog/internal/auth/dto"
)

const (
	credentialsErrorMessage    = "credentials are wrong"
	internalServerErrorMessage = "sorry, try again later"
)

var validate = validator.New()

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

	in := dto.AuthenticateIn(req)
	out, err := h.usecase.Do(ctx, &in)
	if err != nil {
		switch err {
		case dto.ErrUserNotFound, dto.ErrWrongPassword:
			ctx.JSON(http.StatusBadRequest, response{OK: false, Message: credentialsErrorMessage})
		default:
			h.logger.Error("auth: authenticate api error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, response{OK: false, Message: internalServerErrorMessage})
		}
		return
	}

	ctx.JSON(http.StatusOK, response{
		OK:          true,
		AccessToken: out.AccessToken,
	})
}
