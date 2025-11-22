package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/isOdin-l/Assigning-PR-reviewers/configs"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database/postgres"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/httpchi"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg, err := configs.NewConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("Error while initializing config: %v", err.Error()))
		return
	}
	// DB
	db, err := postgres.NewPostgresDB(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while conneciton to database: %v", err.Error()))
	}
	_ = db

	// // repository
	// repository := repository.NewRepo(db)

	// // service
	// service := service.NewService(repository)

	// // handler
	// handler := handler.NewHandler(service)

	// // router
	// router := httpchi.NewRouter(handler)

	// server
	server := httpchi.NewServer("8080", chi.NewRouter())
	go func() {
		if err := server.RunServer(); err != nil {
			slog.Error(fmt.Sprintf("Error while server is running: %v", err.Error()))
			return
		}
	}()
	slog.Info("Server started")

	server.GracefulShutdownServer(ctx)
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}
