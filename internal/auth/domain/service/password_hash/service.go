package password_hash

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/art-es/blog/internal/auth/dto"
)

type Service struct {
	cost int
}

func New() *Service {
	return &Service{
		cost: 10,
	}
}

func (s *Service) Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	return string(hash), err
}

func (s *Service) Validate(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return dto.ErrWrongPassword
	}
	return err
}
