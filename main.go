package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"link_shortener/internal/config"
	"link_shortener/internal/http/handlers"
	"link_shortener/internal/http/middlewares/mwLogger"
	"link_shortener/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("cant open db", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))

	router.Post("/url", handler.Save(storage))
	router.Get("/{alias}", handler.Redirect(storage))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", server.Addr))

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
