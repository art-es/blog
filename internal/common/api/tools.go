package api

import "github.com/art-es/blog/internal/common/log"

type Tools struct {
	Validator          Validator
	ServerErrorHandler ServerErrorHandler
}

func NewTools(logger log.Logger) *Tools {
	return &Tools{
		Validator:          NewValidator(),
		ServerErrorHandler: NewServerErrorHandler(logger),
	}
}
