package routes

import (
	"net/http"
	"path/filepath"
	"url_shortener_v3/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Hm *handler.HandlerManager
}

func InitRoutes(
	h Handlers,
	webDir string,
) *gin.Engine {
	router := gin.Default()

	InitStatic(router, webDir)

	router.GET("/:code", h.Hm.Redirect)

	router.POST("/shorten", h.Hm.Shorten)

	return router
}

func InitStatic(router *gin.Engine, webDir string) {
	router.StaticFile("/", filepath.Join(webDir, "index.html"))

	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
}
