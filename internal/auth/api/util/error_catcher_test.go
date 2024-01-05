package util_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/api/util"
	"github.com/art-es/blog/internal/common/testutil"
)

func TestErrorCatcher(t *testing.T) {
	var (
		loggerBuf  = new(bytes.Buffer)
		handler    = "dummyHandler"
		respWriter = httptest.NewRecorder()
		ctx, _     = gin.CreateTestContext(respWriter)
		err        = errors.New("dummy error")
	)

	c := util.NewErrorCatcher(testutil.NewLoggerWithBuffer(loggerBuf), handler)
	c.Catch(ctx, err)

	logs := strings.Split(loggerBuf.String(), "\n")
	assert.Len(t, logs, 2)

	log := struct {
		Level   string `json:"level"`
		Msg     string `json:"msg"`
		Error   string `json:"error"`
		Handler string `json:"handler"`
	}{}
	_ = json.Unmarshal([]byte(logs[0]), &log)
	assert.Equal(t, "error", log.Level)
	assert.Equal(t, "API internal server error", log.Msg)
	assert.Equal(t, err.Error(), log.Error)
	assert.Equal(t, handler, log.Handler)

	assert.Equal(t, http.StatusInternalServerError, respWriter.Code)
	assert.Equal(t, `{"ok":false,"message":"sorry, try again later"}`, respWriter.Body.String())
}
