package api_test

import (
	"testing"

	"github.com/art-es/blog/internal/common/api"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSetUserID(t *testing.T) {
	const empty = int64(0)
	const notEmpty = int64(1)

	ctx := &gin.Context{}
	assert.Equal(t, empty, api.GetUserID(ctx), "userID is empty")

	api.SetUserID(ctx, notEmpty)
	assert.Equal(t, notEmpty, api.GetUserID(ctx), "userID is not empty")
}
