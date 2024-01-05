package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/art-es/blog/internal/auth/api/activate"
	"github.com/art-es/blog/internal/auth/api/authenticate"
	"github.com/art-es/blog/internal/auth/api/parse_token"
	"github.com/art-es/blog/internal/auth/api/refresh_token"
	"github.com/art-es/blog/internal/auth/api/register"
	auth "github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/domain/service/access_token"
	"github.com/art-es/blog/internal/auth/domain/service/activation"
	"github.com/art-es/blog/internal/auth/domain/service/password_hash"
	"github.com/art-es/blog/internal/auth/infra/databus_kafka"
	"github.com/art-es/blog/internal/auth/infra/repository_pg"
)

func initAuthRoutes(r *gin.Engine, config *Config, logger *zap.Logger, db *sql.DB) {
	repository := repository_pg.New(db)
	databus := databus_kafka.New(config.kafkaURL)
	passwordHashService := password_hash.New()
	activationService := activation.New(logger, databus)
	accessTokenService := access_token.New(config.jwtSecret)

	registerUsecase := auth.NewRegisterUsecase(repository, passwordHashService, activationService)
	r.POST("/v1/auth/register/", register.NewHandler(registerUsecase, logger).Handle)

	activateUsecase := auth.NewActivateUsecase(repository, activationService)
	r.POST("/v1/auth/activate/:code", activate.NewHandler(activateUsecase, logger).Handle)

	authenticateUsecase := auth.NewAuthenticateUsecase(repository.User(), passwordHashService, accessTokenService)
	r.POST("/v1/auth/authenticate", authenticate.NewHandler(authenticateUsecase, logger).Handle)

	refreshTokenUsecase := auth.NewRefreshTokenUsecase(repository.User(), accessTokenService)
	r.POST("/v1/auth/refresh-token", refresh_token.NewHandler(refreshTokenUsecase, logger).Handle)
}

func newParseAccessTokenMiddleware(config *Config, db *sql.DB) gin.HandlerFunc {
	repository := repository_pg.New(db)
	accessTokenService := access_token.New(config.jwtSecret)
	usecase := auth.NewParseTokenUsecase(repository.User(), accessTokenService)

	return parse_token.NewMiddleware(usecase).Handle
}
