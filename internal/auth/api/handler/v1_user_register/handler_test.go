package v1_user_register

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/art-es/blog/internal/auth/api/handler/v1_user_register/mock"
	"github.com/art-es/blog/internal/auth/dto"
	"github.com/art-es/blog/internal/common/api"
	api_mock "github.com/art-es/blog/internal/common/api/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		usecase   = mock.NewMockusecase(ctrl)
		validator = api_mock.NewMockValidator(ctrl)
		router    = gin.New()
		handler   = New(usecase, validator, nil)
	)

	api.BindEndpoint(router, handler)

	var (
		email    = "dummy@example.com"
		password = "1234Qwerty"
	)

	reqBody := bytes.NewBuffer(nil)
	_ = json.NewEncoder(reqBody).Encode(request{Email: email, Password: password})
	validator.EXPECT().Struct(&request{Email: email, Password: password}).Return(nil)
	usecase.EXPECT().Do(gomock.Any(), gomock.Eq(&dto.RegisterIn{Email: email, Password: password})).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(handler.Method(), handler.Endpoint(), reqBody)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"ok":true}`, w.Body.String())
}
