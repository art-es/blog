package main

import (
	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/blog/api/get_article"
	"github.com/art-es/blog/internal/blog/api/get_articles"
	"github.com/art-es/blog/internal/blog/api/get_categories"
)

func initBlogRoutes(
	router *gin.Engine,
	parseTokenMiddleware gin.HandlerFunc,
) {
	// TODO: add implementation
	group := router.Group("/v1/blog", parseTokenMiddleware)
	group.GET("/categories", get_categories.NewHandler().Handle)
	group.GET("/v1/articles", get_articles.NewHandler().Handle)
	group.GET("/v1/articles/:slug", get_article.NewHandler().Handle)
}
