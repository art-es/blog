package apiutil_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/common/apiutil"
)

func TestGetSetUserID(t *testing.T) {
	const empty = int64(0)
	const notEmpty = int64(1)

	ctx := &gin.Context{}
	assert.Equal(t, empty, apiutil.GetUserID(ctx), "userID is empty")

	apiutil.SetUserID(ctx, notEmpty)
	assert.Equal(t, notEmpty, apiutil.GetUserID(ctx), "userID is not empty")
}
