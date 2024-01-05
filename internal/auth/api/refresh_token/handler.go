package refresh_token

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/art-es/blog/internal/auth/api/util"
	"github.com/art-es/blog/internal/auth/dto"
)

const internalServerErrorMessage = "sorry, try again later"

type usecase interface {
	Do(ctx context.Context, in *dto.RefreshTokenIn) (*dto.RefreshTokenOut, error)
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
	accessToken, ok := util.ParseBearerToken(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response{OK: false})
		return
	}

	out, err := h.usecase.Do(ctx, &dto.RefreshTokenIn{AccessToken: accessToken})
	if err != nil {
		switch err {
		case dto.ErrInvalidAccessToken:
			ctx.JSON(http.StatusUnauthorized, response{OK: false})
		default:
			h.logger.Error("auth: refresh token api error", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, response{OK: false, Message: internalServerErrorMessage})
		}
		return
	}

	ctx.JSON(http.StatusOK, response{
		OK:          true,
		AccessToken: out.AccessToken,
	})
}
