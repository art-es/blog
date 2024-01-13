package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/auth/api/handler/v1_token_refresh"
	"github.com/art-es/blog/internal/auth/api/handler/v1_user_activate"
	"github.com/art-es/blog/internal/auth/api/handler/v1_user_authenticate"
	"github.com/art-es/blog/internal/auth/api/handler/v1_user_register"
	"github.com/art-es/blog/internal/auth/api/middleware/parse_token"
	auth "github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/domain/service/access_token"
	"github.com/art-es/blog/internal/auth/domain/service/activation"
	"github.com/art-es/blog/internal/auth/domain/service/password_hash"
	"github.com/art-es/blog/internal/auth/infra/databus_kafka"
	"github.com/art-es/blog/internal/auth/infra/repository_pg"
	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/log"
)

func bindAuthEndpoints(
	r *gin.Engine,
	config *Config,
	logger log.Logger,
	db *sql.DB,
	validator api.Validator,
	serverErrorHandler api.ServerErrorHandler,
) {
	repository := repository_pg.New(db)
	databus := databus_kafka.New(config.kafkaURL)
	passwordHashService := password_hash.New()
	activationService := activation.New(logger, databus)
	accessTokenService := access_token.New(config.jwtSecret)

	api.BindEndpoints(r,
		v1_user_register.New(
			auth.NewRegisterUsecase(repository, passwordHashService, activationService),
			validator,
			serverErrorHandler,
		),
		v1_user_activate.New(
			auth.NewActivateUsecase(repository, activationService),
			serverErrorHandler,
		),
		v1_user_authenticate.New(
			auth.NewAuthenticateUsecase(repository.User(), passwordHashService, accessTokenService),
			validator,
			serverErrorHandler,
		),
		v1_token_refresh.New(
			auth.NewRefreshTokenUsecase(repository.User(), accessTokenService),
			serverErrorHandler,
		),
	)
}

func newParseAccessTokenMiddleware(config *Config, db *sql.DB) gin.HandlerFunc {
	repository := repository_pg.New(db)
	accessTokenService := access_token.New(config.jwtSecret)
	usecase := auth.NewParseTokenUsecase(repository.User(), accessTokenService)

	return parse_token.New(usecase).Handle
}
