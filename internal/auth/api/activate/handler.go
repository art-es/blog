package activate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/art-es/blog/internal/auth/dto"
)

const (
	invalidCodeErrorMessage    = "invalid code"
	internalServerErrorMessage = "sorry, try again later"
)

type usecase interface {
	Do(ctx context.Context, in *dto.ActivateIn) error
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
	code := ctx.Param("code")
	if code == "" || !validUUID(code) {
		ctx.JSON(http.StatusBadRequest, response{OK: false, Message: invalidCodeErrorMessage})
		return
	}

	err := h.usecase.Do(ctx, &dto.ActivateIn{Code: code})
	if err != nil {
		switch err {
		case dto.ErrActivationCodeNotFound, dto.ErrUserNotFound:
			ctx.JSON(http.StatusNotFound, response{OK: false})
		default:
			h.logger.Error("auth: activate api error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, response{OK: false, Message: internalServerErrorMessage})
		}
		return
	}

	ctx.JSON(http.StatusOK, response{OK: true})
}

func validUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
