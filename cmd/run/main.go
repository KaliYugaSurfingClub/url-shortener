package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"shortener/internal/config"
	"shortener/internal/core/generator"
	"shortener/internal/core/model"
	"shortener/internal/core/services/adViewer"
	"shortener/internal/core/services/linkManager"
	"shortener/internal/core/services/linkShortener"
	"shortener/internal/storage/postgres"
	"shortener/internal/storage/postgres/repository"
	"shortener/internal/transport/rest/handler"
	"shortener/internal/transport/rest/handler/completeAdHandler"
	"shortener/internal/transport/rest/handler/getLinkClicksHandler"
	"shortener/internal/transport/rest/handler/getUserLinksHandler"
	"shortener/internal/transport/rest/handler/openShortenedHandler"
	"shortener/internal/transport/rest/handler/shortLinkHandler"
	"shortener/internal/transport/rest/mw"
	"time"
)

type FakePayer struct{}

func (f FakePayer) Pay(ctx context.Context, clickId int64) error {
	return nil
}

type FakeAdProvider struct{}

func (f FakeAdProvider) GetAdByMetadata(ctx context.Context, metadata model.ClickMetadata) (int64, error) {
	return 1, nil
}

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, cancel, err := postgres.NewPgxPool(cfg.PostgresURL)
	if err != nil {
		log.Error("unable to connect to postgres", mw.ErrAttr(err))
		cancel()
	}

	defer cancel()

	repo := repository.New(db)

	aliasGenerator := generator.New([]rune("abcdefgr"), 4)
	aliasManager, err := linkShortener.New(repo, aliasGenerator, 1)
	manager := linkManager.New(repo)

	adViewManager := adViewer.New(repo, &FakePayer{}, &FakeAdProvider{})
	onCompleteErrs := adViewManager.OnCompleteErrs()

	go func() {
		for err := range onCompleteErrs {
			log.Error(err.Error())
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

	videoTemplate, err := template.ParseFiles("C:\\Users\\leono\\Desktop\\prog\\go\\shortener\\adPages\\ADVideo.html")
	if err != nil {
		log.Error(err.Error())
	}

	fileTemplate, err := template.ParseFiles("C:\\Users\\leono\\Desktop\\prog\\go\\shortener\\adPages\\ADFile.html")
	if err != nil {
		log.Error(err.Error())
	}

	templates := map[model.AdType]*template.Template{
		model.AdTypeVideo: videoTemplate,
		model.AdTypeFile:  fileTemplate,
	}

	r.Get("/{alias}", openShortenedHandler.New(adViewManager, templates, "http://localhost:8080/static/video", "http://localhost:8080").Handler)
	r.Post("/{id}", completeAdHandler.New(adViewManager).Handler)
	r.Get("/static/video", handler.StreamVideoHandler)

	server := &http.Server{
		Addr:              "0.0.0.0:8080",
		Handler:           r,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Unable to start server: ", err)
		os.Exit(1) //todo
	}
}
