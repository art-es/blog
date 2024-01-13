package testutil

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/common/api"
)

func NewGinRouter(handler api.EndpointHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.BindEndpoint(router, handler)
	return router
}

func ReadCloserFromJSON(obj any) io.ReadCloser {
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(obj)
	return io.NopCloser(buf)
}
