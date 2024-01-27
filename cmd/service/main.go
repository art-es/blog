package main

import (
	"database/sql"
	"fmt"

	"github.com/art-es/blog/internal/auth/infra/databus_kafka"

	"github.com/art-es/blog/internal/auth/domain/service/access_token"

	"github.com/art-es/blog/internal/common/validation"

	"github.com/art-es/blog/internal/auth/domain/service/activation"
	"github.com/art-es/blog/internal/auth/domain/service/password_hash"
	"github.com/art-es/blog/internal/auth/infra/repository_pg"

	"github.com/art-es/blog/internal/auth/api/endpoint/v1_user_register"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/cmd/service/config"
	"github.com/art-es/blog/internal/auth/api/endpoint/v1_access_token_refresh"
	"github.com/art-es/blog/internal/auth/api/endpoint/v1_user_activate"
	"github.com/art-es/blog/internal/auth/api/endpoint/v1_user_authenticate"
	auth "github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/log"
)

func main() {
	panic(run())
}

func run() error {
	conf := config.Parse()
	logger := log.New()

	db, err := sql.Open("postgres", conf.PGConnect)
	if err != nil {
		return fmt.Errorf("open database error: %w", err)
	}

	router := gin.Default()
	bindAuthEndpoints(router, conf, logger, db)

	err = router.Run(conf.ServiceURL)
	return fmt.Errorf("running router error: %w", err)
}

func bindAuthEndpoints(router *gin.Engine, conf *config.Config, logger log.Logger, db *sql.DB) {
	repository := repository_pg.New(db)
	passwordHashService := password_hash.New()
	activationService := activation.New(logger, databus_kafka.New(conf.KafkaURL))
	accessTokenService := access_token.New(conf.AccessTokenSecret)

	validator := validation.NewValidator()
	serverErrorHandlerFactory := api.NewServerErrorHandlerFactory(logger)

	v1_user_register.Bind(
		router,
		auth.NewUserRegisterCase(repository, passwordHashService, activationService),
		validator,
		serverErrorHandlerFactory,
	)
	v1_user_activate.Bind(
		router,
		auth.NewUserActivateCase(repository, activationService),
		validator,
		serverErrorHandlerFactory,
	)
	v1_user_authenticate.Bind(
		router,
		auth.NewUserAuthenticateCase(repository.User(), passwordHashService, accessTokenService),
		validator,
		serverErrorHandlerFactory,
	)
	v1_access_token_refresh.Bind(
		router,
		auth.NewAccessTokenRefreshCase(repository.User(), accessTokenService),
		serverErrorHandlerFactory,
	)
}
