package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/common/api"
	"github.com/art-es/blog/internal/common/log"
)

func main() {
	config := getConfig()
	logger := log.New()
	pgConn, err := sql.Open("postgres", config.pgConnect)
	if err != nil {
		panic(err)
	}

	router := gin.New()

	var (
		validator                  = api.NewValidator()
		serverErrorHandler         = api.NewServerErrorHandler(logger)
		parseAccessTokenMiddleware = newParseAccessTokenMiddleware(config, pgConn)
	)

	bindAuthEndpoints(router, config, logger, pgConn, validator, serverErrorHandler)
	bindBlogEndpoints(router, parseAccessTokenMiddleware)

	if err = router.Run(config.serviceUrl); err != nil {
		panic(err)
	}
}
