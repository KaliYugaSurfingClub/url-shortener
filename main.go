package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"link_shortener/internal/config"
	"link_shortener/internal/storage/sqlite/alias"
	"link_shortener/internal/storage/sqlite/visit"
	"link_shortener/internal/transport/handlers/redirect"
	"link_shortener/internal/transport/handlers/save"
	"link_shortener/internal/transport/middlewares/mwLogger"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	//todo remove literal
	storage, err := alias.New(cfg.StoragePath)
	if err != nil {
		log.Error("cant open db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	visitStorage, err := visit.New(cfg.StoragePath, 1*time.Minute)
	if err != nil {
		log.Error("cant open db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log.Info("open db")

	al, err := visitStorage.GetUnusedAliases()
	fmt.Println(err)
	fmt.Println(al)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))

	router.Post("/url", save.New(storage))
	router.Get("/{alias}", redirect.New(storage))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", server.Addr))

	//go cleaner.Start(log, storage, 2*time.Minute, time.Minute)

	if err := server.ListenAndServe(); err != nil {
		log.Error("cant start server", slog.String("error", err.Error()))
	}

	log.Error("server stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: false,
			},
		))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			},
		))
	}

	return logger
}
