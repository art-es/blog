package password_hash

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/art-es/blog/internal/auth/dto"
)

func TestService_happyPath(t *testing.T) {
	const password = "qwer1234!"
	s := New()

	hash, err := s.Generate(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = s.Validate(hash, password)
	assert.NoError(t, err)
}

func TestService_wrongPassword(t *testing.T) {
	const password = "qwer1234!"
	s := New()

	hash, err := s.Generate(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = s.Validate(hash, "fake-password")
	assert.ErrorIs(t, err, dto.ErrIncorrectPassword)
}
