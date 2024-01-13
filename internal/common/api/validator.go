//go:generate mockgen -source=validator.go -destination=mock/validator.go -package=mock
package api

import "github.com/go-playground/validator/v10"

type Validator interface {
	Struct(s interface{}) error
}

func NewValidator() *validator.Validate {
	return validator.New()
}
