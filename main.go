package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"link_shortener/internal/config"
	"link_shortener/internal/storage/sqlite/urlManager"
	"log"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = urlManager.New(db)
	if err != nil {
		log.Fatal(err)
	}

	//log := setupLogger(cfg.Env)
	//
	////todo remove literal
	//storage, err := aliasStorage.New(cfg.StoragePath)
	//if err != nil {
	//	log.Error("cant open db", slog.String("error", err.Error()))
	//	os.Exit(1)
	//}
	//
	//visitStorage, err := visit.New(cfg.StoragePath, 1*time.Minute)
	//if err != nil {
	//	log.Error("cant open db", slog.String("error", err.Error()))
	//	os.Exit(1)
	//}
	//
	//log.Info("open db")
	//
	//al, err := visitStorage.GetUnusedAliases()
	//fmt.Println(err)
	//fmt.Println(al)
	//
	//router := chi.NewRouter()
	//
	//router.Use(middleware.RequestID)
	//router.Use(middleware.Recoverer)
	//router.Use(middleware.URLFormat)
	//router.Use(mwLogger.New(log))
	//
	//router.Post("/url", save.New(storage))
	//router.Get("/{alias}", redirect.New(storage))
	//
	//server := &http.Server{
	//	Addr:         cfg.Address,
	//	Handler:      router,
	//	ReadTimeout:  cfg.Timeout,
	//	WriteTimeout: cfg.Timeout,
	//	IdleTimeout:  cfg.IdleTimeout,
	//}
	//
	//log.Info("starting server", slog.String("address", server.Addr))
	//
	////go cleaner.Start(log, storage, 2*time.Minute, time.Minute)
	//
	//if err := server.ListenAndServe(); err != nil {
	//	log.Error("cant start server", slog.String("error", err.Error()))
	//}
	//
	//log.Error("server stopped")
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
