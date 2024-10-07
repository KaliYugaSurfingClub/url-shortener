package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"shortener/internal/core/generator"
	"shortener/internal/core/model"
	"shortener/internal/core/services/adViewer"
	"shortener/internal/core/services/linkManager"
	"shortener/internal/core/services/linkShortener"
	"shortener/internal/storage/postgres/clickRepo"
	"shortener/internal/storage/postgres/linkRepo"
	"shortener/internal/storage/transaction"
	"shortener/internal/transport/rest/handler"
	"shortener/internal/transport/rest/handler/getLinkClicksHandler"
	"shortener/internal/transport/rest/handler/getUserLinksHandler"
	"shortener/internal/transport/rest/handler/openShortenedHandler"
	"shortener/internal/transport/rest/handler/shortLinkHandler"
	"shortener/internal/transport/rest/mw"
	"time"
)

type temporaryNotifier struct{}

func (n *temporaryNotifier) NotifyOpen(context.Context, *model.Link, int64)      {}
func (n *temporaryNotifier) NotifyClosed(context.Context, *model.Link, int64)    {}
func (n *temporaryNotifier) NotifyCompleted(context.Context, *model.Link, int64) {}

type tempPayer struct{}

func (t *tempPayer) Pay(context.Context, int64) error { return nil }

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
	transactor := transaction.NewTransactor(db)

	aliasGenerator := generator.New([]rune("abcdefgr"), 4)
	aliasManager, err := linkShortener.New(linkStore, aliasGenerator, 10)
	manager := linkManager.New(linkStore, clickStore)

	adViewManager := adViewer.New(linkStore, clickStore, &tempPayer{}, &temporaryNotifier{}, transactor)

	jwtOpt := mw.JwtOptions{
		UserIdKey:  "id",
		CookieName: "user_id",
		Secret:     []byte("sasha"),
	}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mw.NewLogger(log))

	r.Route("/users", func(r chi.Router) {
		r.Get("/login", handler.Login(jwtOpt, 24*time.Hour))
	})

	r.Route("/links", func(r chi.Router) {
		r.Use(mw.CheckAuth(jwtOpt))
		r.Post("/", shortLinkHandler.New(aliasManager, 255, 255, 255).Handler)
		r.Get("/", getUserLinksHandler.New(manager).Handler)
		r.Get("/{id}/clicks", getLinkClicksHandler.New(manager).Handler)
	})

	t, _ := template.ParseFiles("C:\\Users\\leono\\Desktop\\prog\\go\\shortener\\adPages\\AD.html")

	r.Get("/{alias}", openShortenedHandler.New(adViewManager, "/static/video", t).Handler)

	r.Get("/static/video", handler.StreamVideoHandler) //todo maybe /static/video may returns random video

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Unable to start server: ", err)
		os.Exit(1) //todo
	}
}
