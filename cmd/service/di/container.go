package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/cmd/service/config"
	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/log"
)

type Middlewares struct {
	ParseAccessToken gin.HandlerFunc
	Authenticated    gin.HandlerFunc
}

type Container struct {
	Config      *config.Config
	Logger      log.Logger
	DB          *sql.DB
	APITools    *api.Tools
	Middlewares *Middlewares
	Auth        *AuthContainer
}

func New(conf *config.Config, logger log.Logger, db *sql.DB) *Container {
	c := Container{
		Config:   conf,
		Logger:   logger,
		DB:       db,
		APITools: api.NewTools(logger),
	}

	bindAuth(&c)

	return &c
}
