package access_token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/domain/service/access_token"
)

func TestService_happyPath(t *testing.T) {
	const userID = int64(1)

	service := access_token.New([]byte("secret"))

	object := service.NewObject(userID)
	assertObject(t, object, userID)

	service.Refresh(object)
	assertObject(t, object, userID)

	token, err := service.Sign(object)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parserResult, err := service.Parse(token)
	assert.NoError(t, err)
	assertEqualObjects(t, object, parserResult)
}

func TestService_ParseAndValidate_expired(t *testing.T) {
	const userID = int64(1)

	service := access_token.New([]byte("secret"))

	object := service.NewObject(userID)
	object.ExpirationTime = time.Now().Add(-3 * time.Hour)

	token, err := service.Sign(object)
	assert.NoError(t, err)

	parsedObject, err := service.ParseAndValidate(token)
	assert.EqualError(t, err, "parse string with secret error: token has invalid claims: token is expired")
	assert.Nil(t, parsedObject)
}

func assertObject(t *testing.T, object *domain.AccessTokenObject, expUserID int64) {
	assert.NotEqual(t, time.Time{}, object.ExpirationTime)
	assert.Equal(t, object.NotBefore, object.IssuedAt)
	assert.Equal(t, 12*time.Hour, object.ExpirationTime.Sub(object.IssuedAt))
	assert.Equal(t, []string{"auth"}, object.Audience)
	assert.Equal(t, "art-es", object.Issuer)
	assert.Equal(t, expUserID, object.UserID)
}

func assertEqualObjects(t *testing.T, expected, actual *domain.AccessTokenObject) {
	assertTimes := func(expected, actual time.Time) {
		const format = time.RFC3339
		assert.Equal(t, expected.Format(format), actual.Format(format))
	}

	assertTimes(expected.ExpirationTime, actual.ExpirationTime)
	assertTimes(expected.NotBefore, actual.NotBefore)
	assertTimes(expected.IssuedAt, actual.IssuedAt)
	assert.Equal(t, expected.Audience, actual.Audience)
	assert.Equal(t, expected.Issuer, actual.Issuer)
	assert.Equal(t, expected.UserID, actual.UserID)
}
