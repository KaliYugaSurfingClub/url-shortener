package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"os"
	"shortener/internal/core/generator"
	"shortener/internal/core/services/clickManager"
	"shortener/internal/core/services/linkManager"
	"shortener/internal/core/services/linkShortener"
	"shortener/internal/storage/postgres/clickRepo"
	"shortener/internal/storage/postgres/linkRepo"
	"shortener/internal/transport/rest/handler"
	"shortener/internal/transport/rest/handler/getLinkClicksHandler"
	"shortener/internal/transport/rest/handler/getUserLinksHandler"
	"shortener/internal/transport/rest/handler/shortLinkHandler"
	"shortener/internal/transport/rest/mw"
	"time"
)

type FakeUserStore struct{}

func (u *FakeUserStore) AddToBalance(ctx context.Context, id int64, payment int) error {
	fmt.Println("balance increasing")
	return nil
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	dbURL := "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"

	poolCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		fmt.Println("Unable to parse DATABASE_URL: ", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		fmt.Println("Unable to create connection pool: ", err)
	}

	defer db.Close()

	linkStore := linkRepo.New(db)
	clickStore := clickRepo.New(db)
	//userStore := &FakeUserStore{}
	//transactor := transaction.NewTransactor(db)

	aliasGenerator := generator.New([]rune("abcdefgr"), 4)
	aliasManager, err := linkShortener.New(linkStore, aliasGenerator, 10)
	linkM := linkManager.New(linkStore)
	clickM := clickManager.New(clickStore)

	//adViewManager := adViewManager.New(linkStore, clickStore, userStore, transactor)

	jwtOpt := mw.JwtOptions{
		UserIdKey:  "id",
		CookieName: "user_id",
		Secret:     []byte("sasha"),
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mw.NewLogger(log))

	r.Route("/user", func(r chi.Router) {
		r.Get("/login", handler.Login(jwtOpt, 24*time.Hour))
	})

	r.Route("/link", func(r chi.Router) {
		r.Use(mw.CheckAuth(jwtOpt))
		r.Post("/", shortLinkHandler.New(aliasManager, 255, 255, 255).Handler)
		r.Get("/", getUserLinksHandler.New(linkM, 10).Handler)
		r.Get("/{id}/clicks", getLinkClicksHandler.New(clickM, 10).Handler)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Unable to start server: ", err)
		os.Exit(1) //todo
	}
}
