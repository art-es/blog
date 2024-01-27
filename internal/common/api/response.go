package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error   *Error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func RequestValidationFailedResponse(ctx *gin.Context, err error) {
	var message string
	if err != nil {
		message = err.Error()
	}

	ctx.JSON(http.StatusBadRequest, &ErrorResponse{
		Error: &Error{
			Code: 1001,
			Name: "Request validation failed",
		},
		Message: message,
	})
}

func UnauthorizedResponse(ctx *gin.Context) {
	const message = "Please try to sign in again."
	ctx.JSON(http.StatusUnauthorized, &ErrorResponse{Message: message})
}

func InternalServerErrorResponse(ctx *gin.Context) {
	const message = "Something went wrong, please try again later."
	ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Message: message})
}
