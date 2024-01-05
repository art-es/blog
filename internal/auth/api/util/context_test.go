package util_test

import (
	"testing"

	"github.com/art-es/blog/internal/auth/api/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSetUserID(t *testing.T) {
	const empty = int64(0)
	const notEmpty = int64(1)

	ctx := &gin.Context{}
	assert.Equal(t, empty, util.GetUserID(ctx))
	util.SetUserID(ctx, notEmpty)
	assert.Equal(t, notEmpty, util.GetUserID(ctx))
}
