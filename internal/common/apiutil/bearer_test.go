package apiutil_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/common/apiutil"
)

func TestParseBearerToken(t *testing.T) {
	tests := []struct {
		headerValue string
		expToken    string
		expExists   bool
	}{
		{
			headerValue: "Bearer foo",
			expToken:    "foo",
			expExists:   true,
		},
		{
			headerValue: "",
			expExists:   false,
		},
		{
			headerValue: "Bearer",
			expExists:   false,
		},
		{
			headerValue: "Bearer ",
			expExists:   false,
		},
		{
			headerValue: "Something foo",
			expExists:   false,
		},
		{
			headerValue: "Bearer foo bar",
			expExists:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.headerValue, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://127.0.0.1", nil)
			req.Header.Set("Authorization", tt.headerValue)
			ctx := &gin.Context{Request: req}

			token, exists := apiutil.ParseBearerToken(ctx)
			assert.Equal(t, tt.expToken, token)
			assert.Equal(t, tt.expExists, exists)
		})
	}
}
