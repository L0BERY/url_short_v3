package routes

import (
	"url_shortener_v3/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Hm handler.HandlerManager
}

func InitRoutes(
	h Handlers,
) *gin.Engine {
	router := gin.Default()

	router.GET("/:code", h.Hm.Redirect)

	router.POST("/shorten", h.Hm.Shorten)

	return router
}

// func InitStatic(router *gin.Engine, webDir string) {
// 	static := router.Group("/", func(ctx *gin.Context) {
// 		ctx.Header("Cache-control", "no-cache")
// 	})

// 	static.Static()
// }
