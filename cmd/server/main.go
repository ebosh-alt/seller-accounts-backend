package main

//go:generate sh -c "cd ../.. && swag init -g main.go -o docs -d cmd/server,internal/delivery/http/server,internal/entities,internal/usecase,internal/repository"

import (
	"context"
	"os"
	"os/signal"
	"sellers-accounts-backend/internal/service"
	"syscall"
	"time"

	"sellers-accounts-backend/config"
	"sellers-accounts-backend/internal/delivery/http/server"
	"sellers-accounts-backend/internal/delivery/http/server/handlers"
	"sellers-accounts-backend/internal/delivery/http/server/middleware"
	"sellers-accounts-backend/internal/repository"
	"sellers-accounts-backend/internal/usecase"
	"sellers-accounts-backend/pkg/logger"

	"github.com/gin-gonic/gin"

	_ "sellers-accounts-backend/docs"
)

// @title Sellers Accounts API
// @version 1.0.0
// @description HTTP API for sellers accounts.
// @BasePath /api
// @schemes https http
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.New("debug")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	repo, err := repository.New("postgres", log, cfg, ctx)
	if err != nil {
		log.Fatalw("failed to init repository", "error", err)
	}
	if err := repo.OnStart(ctx); err != nil {
		log.Fatalw("failed to start repository", "error", err)
	}
	defer func() {
		if err := repo.OnStop(ctx); err != nil {
			log.Errorw("failed to stop repository", "error", err)
		}
	}()

	botLinkService := cache.NewCache(repo, time.Duration(cfg.Cache.TTLSeconds))
	uc := usecase.New(cfg, log, ctx, repo, botLinkService)

	engine := gin.Default()
	engine.Use(gin.Recovery())

	mdl := middleware.New(cfg, log, engine)
	handlers := handler.New(log, engine, uc)
	httpServer := server.New(log, cfg, engine, handlers, *mdl)

	if err := httpServer.OnStart(); err != nil {
		log.Fatalw("failed to start server", "error", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.OnStop(shutdownCtx); err != nil {
		log.Errorw("failed to stop server", "error", err)
	}
}
