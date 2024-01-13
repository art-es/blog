//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package v1_token_refresh

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/apiutil"
)

type usecase interface {
	Do(ctx context.Context, in *dto.RefreshTokenIn) (*dto.RefreshTokenOut, error)
}

type response struct {
	OK          bool   `json:"ok"`
	AccessToken string `json:"accessToken,omitempty"`
	Message     string `json:"message,omitempty"`
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
	return http.MethodPost
}

func (h *Handler) Endpoint() string {
	return "/v1/auth/token/refresh"
}

func (h *Handler) Handle(ctx *gin.Context) {
	token, ok := apiutil.ParseBearerToken(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response{})
		return
	}

	in := dto.RefreshTokenIn{AccessToken: token}
	out, err := h.usecase.Do(ctx, &in)
	if err != nil {
		switch err {
		case dto.ErrInvalidAccessToken:
			ctx.JSON(http.StatusUnauthorized, response{})
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
