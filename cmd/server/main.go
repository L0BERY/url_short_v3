package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_shortener_v3/internal/config"
	"url_shortener_v3/internal/database"
	"url_shortener_v3/internal/handler"
	"url_shortener_v3/internal/repository"
	"url_shortener_v3/internal/routes"
	"url_shortener_v3/internal/service"
)

func main() {
	checkHealth := flag.Bool("check-health", false, "опросить запущенный сервер и выйти")
	flag.Parse()

	cfg := config.LoadConfig()

	if *checkHealth {
		os.Exit(probeHealth(cfg.ServerAddress))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.ConnectPostgres(ctx, cfg.DSN)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	hm := handler.NewHandlerManager(service.NewService(repository.NewRepository(pool)))
	router := routes.InitRoutes(routes.Handlers{Hm: hm}, "web")

	httpSrv := &http.Server{
		Addr:              cfg.ServerAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()
	log.Printf("listening on %s", cfg.ServerAddress)

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func probeHealth(addr string) (exitCode int) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Printf("health: некорректный адрес %q: %v", addr, err)
		return 1
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}

	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get("http://" + net.JoinHostPort(host, port) + "/")
	if err != nil {
		log.Printf("health: %v", err)
		return 1
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %v", closeErr)
			if exitCode == 0 {
				exitCode = 1
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("health: HTTP %d", resp.StatusCode)
		return 1
	}

	return 0
}
