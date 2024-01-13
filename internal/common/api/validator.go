//go:generate mockgen -source=validator.go -destination=mock/validator.go -package=mock
package api

import base "github.com/go-playground/validator/v10"

var _ Validator = (*base.Validate)(nil)

type Validator interface {
	Struct(s interface{}) error
	Var(field interface{}, tag string) error
}

type validator struct {
	*base.Validate
}

func NewValidator() *validator {
	return &validator{
		base.New(),
	}
}
