package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/common/api"
)

func BusyEmailResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
		Error: &api.Error{
			Code: 2001,
			Name: "Busy email",
		},
		Message: "User with this email already exists.",
	})
}

func ExpiredUserActivationCodeResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
		Error: &api.Error{
			Code: 2002,
			Name: "Expired user activation code",
		},
		Message: "Activate code is expired. Please pass the registration process again.",
	})
}

func IncorrectUserCredentialsResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, &api.ErrorResponse{
		Error: &api.Error{
			Code: 2003,
			Name: "Incorrect user credentials",
		},
		Message: "Email or password is incorrect.",
	})
}
