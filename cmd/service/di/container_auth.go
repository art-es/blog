package di

import (
	"github.com/art-es/blog/internal/auth/api/middleware/authenticated"
	"github.com/art-es/blog/internal/auth/api/middleware/parse_token"
	auth "github.com/art-es/blog/internal/auth/domain"
	"github.com/art-es/blog/internal/auth/domain/service/access_token"
	"github.com/art-es/blog/internal/auth/domain/service/activation"
	"github.com/art-es/blog/internal/auth/domain/service/password_hash"
	"github.com/art-es/blog/internal/auth/infra/databus_kafka"
	"github.com/art-es/blog/internal/auth/infra/repository_pg"
)

type AuthContainer struct {
	Repository          *repository_pg.Repository
	Databus             *databus_kafka.Client
	PasswordHashService *password_hash.Service
	ActivationService   *activation.Service
	AccessTokenService  *access_token.Service
}

func bindAuth(c *Container) {
	var (
		repository          = repository_pg.New(c.DB)
		databus             = databus_kafka.New(c.Config.KafkaURL)
		passwordHashService = password_hash.New()
		activationService   = activation.New(c.Logger, databus)
		accessTokenService  = access_token.New(c.Config.AccessTokenSecret)
	)

	c.Auth = &AuthContainer{
		Repository:          repository,
		Databus:             databus,
		PasswordHashService: passwordHashService,
		ActivationService:   activationService,
		AccessTokenService:  accessTokenService,
	}

	c.Middlewares.ParseAccessToken = parse_token.New(auth.NewParseTokenUsecase(repository.User(), accessTokenService)).Handle
	c.Middlewares.Authenticated = authenticated.New().Handle
}
