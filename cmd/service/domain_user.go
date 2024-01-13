package main

import (
	"github.com/gin-gonic/gin"

	"github.com/art-es/blog/internal/blog/api/get_article"
	"github.com/art-es/blog/internal/blog/api/get_articles"
	"github.com/art-es/blog/internal/blog/api/get_categories"
)

func bindBlogEndpoints(
	r *gin.Engine,
	parseTokenMiddleware gin.HandlerFunc,
) {
	// TODO: add implementation
	group := r.Group("/v1/blog", parseTokenMiddleware)
	group.GET("/categories", get_categories.NewHandler().Handle)
	group.GET("/articles", get_articles.NewHandler().Handle)
	group.GET("/articles/:slug", get_article.NewHandler().Handle)
}
