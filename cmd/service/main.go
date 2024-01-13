package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/cmd/service/config"
	"github.com/art-es/blog/cmd/service/di"
	"github.com/art-es/blog/internal/auth/api/handler/v1_token_refresh"
	"github.com/art-es/blog/internal/auth/api/handler/v1_user_activate"
	"github.com/art-es/blog/internal/auth/api/handler/v1_user_authenticate"
	"github.com/art-es/blog/internal/auth/api/handler/v1_user_register"
	authDomain "github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/blog/api/get_article"
	"github.com/art-es/blog/internal/blog/api/get_articles"
	"github.com/art-es/blog/internal/blog/api/get_categories"
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

	dic := di.New(conf, logger, db)
	router := gin.Default()

	bindAuthEndpoints(router, dic)
	bindBlogEndpoints(router, dic)

	err = router.Run(conf.ServiceURL)
	return fmt.Errorf("running router error: %w", err)
}

func bindAuthEndpoints(r *gin.Engine, dic *di.Container) {
	authDI := dic.Auth
	apiTools := dic.APITools

	api.BindEndpoints(r,
		v1_user_register.New(
			authDomain.NewRegisterUsecase(authDI.Repository, authDI.PasswordHashService, authDI.ActivationService),
			apiTools.Validator,
			apiTools.ServerErrorHandler,
		),
		v1_user_activate.New(
			authDomain.NewActivateUsecase(authDI.Repository, authDI.ActivationService),
			apiTools.Validator,
			apiTools.ServerErrorHandler,
		),
		v1_user_authenticate.New(
			authDomain.NewAuthenticateUsecase(authDI.Repository.User(), authDI.PasswordHashService, authDI.AccessTokenService),
			apiTools.Validator,
			apiTools.ServerErrorHandler,
		),
		v1_token_refresh.New(
			authDomain.NewRefreshTokenUsecase(authDI.Repository.User(), authDI.AccessTokenService),
			apiTools.ServerErrorHandler,
		),
	)
}

func bindBlogEndpoints(r *gin.Engine, dic *di.Container) {
	// TODO: add implementation
	rg := r.Group("/v1/blog", dic.Middlewares.ParseAccessToken)
	rg.GET("/categories", get_categories.NewHandler().Handle)
	rg.GET("/articles", get_articles.NewHandler().Handle)
	rg.GET("/articles/:slug", get_article.NewHandler().Handle)
}
