package parse_token_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/art-es/blog/internal/common/api"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/api/middleware/parse_token"
	"github.com/art-es/blog/internal/auth/api/middleware/parse_token/mock"
	"github.com/art-es/blog/internal/auth/dto"
)

func Test(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usecase := mock.NewMockusecase(ctrl)

	gin.SetMode(gin.TestMode)
	middleware := parse_token.New(usecase).Handle

	tests := []struct {
		name    string
		setup   func(r *http.Request)
		handler func(c *gin.Context)
	}{
		{
			name: "token parsed",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer foo")

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.AccessTokenParseIn{AccessToken: "foo"})).
					Return(&dto.ParseTokenOut{UserID: 1}, nil)
			},
			handler: func(c *gin.Context) {
				userID := api.GetUserID(c)
				assert.Equal(t, int64(1), userID)
			},
		},
		{
			name: "token not specified",
			setup: func(r *http.Request) {
				r.Header.Del("Authorization")
			},
			handler: func(c *gin.Context) {
				userID := api.GetUserID(c)
				assert.Equal(t, int64(0), userID)
			},
		},
		{
			name: "token parsing error",
			setup: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer foo")

				usecase.EXPECT().
					Do(gomock.Any(), gomock.Eq(&dto.AccessTokenParseIn{AccessToken: "foo"})).
					Return(nil, errors.New("dummy error"))
			},
			handler: func(c *gin.Context) {
				userID := api.GetUserID(c)
				assert.Equal(t, int64(0), userID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const method = http.MethodGet
			const path = "/"

			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, path, nil)
			tt.setup(r)

			router := gin.New()
			router.Use(middleware)
			router.Handle(method, path, tt.handler)
			router.ServeHTTP(w, r)
		})
	}
}
