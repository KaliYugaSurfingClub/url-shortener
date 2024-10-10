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

type TemporaryNotifier struct{}

func (t *TemporaryNotifier) Notify(_ context.Context, eventType model.AdStatus, link *model.Link, click *model.Click) {
	fmt.Printf("%s clickId - %+v link - %+v", eventType, click, link)
}

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

	adViewManager := adViewer.New(linkStore, clickStore, &tempPayer{}, &TemporaryNotifier{}, transactor, 16, 16)
	onCompleteErrs := adViewManager.OnCompleteErrs()

	go func() {
		for err := range onCompleteErrs {
			fmt.Println(err)
		}
	}()

	go adViewManager.StartCleaningExpiredSessions(1*time.Minute, 1*time.Minute, 100)
	cleanerErrs := adViewManager.CleanerErrs()

	go func() {
		for err := range cleanerErrs {
			fmt.Println(err)
		}
	}()

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

	r.Get("/static/video", handler.StreamVideoHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Unable to start server: ", err)
		os.Exit(1) //todo
	}
}
