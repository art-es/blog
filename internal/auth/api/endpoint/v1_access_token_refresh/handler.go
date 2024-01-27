package v1_access_token_refresh

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
)

type response struct {
	AccessToken string `json:"accessToken,omitempty"`
}

type handler struct {
	accessTokenRefreshCase accessTokenRefreshCase
	serverErrorHandler     api.ServerErrorHandler
}

func (h *handler) handle(ctx *gin.Context) {
	smellyAccessToken := ctx.GetHeader("X-Access-Token")
	if smellyAccessToken == "" {
		api.UnauthorizedResponse(ctx)
		return
	}

	freshAccessToken, err := h.useCase(ctx, smellyAccessToken)
	if err != nil {
		switch err {
		case dto.ErrInvalidAccessToken:
			api.UnauthorizedResponse(ctx)
		default:
			h.serverErrorHandler.Handle(ctx, err)
		}
		return
	}

	okResponse(ctx, freshAccessToken)
}

func (h *handler) useCase(ctx *gin.Context, accessToken string) (string, error) {
	in := dto.AccessTokenRefreshIn{
		AccessToken: accessToken,
	}

	out, err := h.accessTokenRefreshCase.Use(ctx, &in)
	if err != nil {
		return "", err
	}

	return out.AccessToken, nil
}

func okResponse(ctx *gin.Context, accessToken string) {
	ctx.JSON(http.StatusOK, &response{AccessToken: accessToken})
}
