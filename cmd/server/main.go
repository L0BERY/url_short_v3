package main

import (
	"context"
	"url_shortener_v3/internal/config"
	"url_shortener_v3/internal/database"
	"url_shortener_v3/internal/handler"
	"url_shortener_v3/internal/repository"
	"url_shortener_v3/internal/routes"
	"url_shortener_v3/internal/service"
)

func main() {

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := database.ConnectPostgres(ctx, cfg.DSN)
	if err != nil {
		return
	}

	repo := repository.NewRepository(conn)

	srv := service.NewService(repo)

	handl := handler.NewHandlerManager(srv)

	routes.InitRoutes(routes.Handlers{Hm: handl}, "web")
}
