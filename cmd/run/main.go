package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"shortener/internal/config"
	"shortener/internal/core/generator"
	"shortener/internal/core/model"
	"shortener/internal/core/services/adViewer"
	"shortener/internal/core/services/linkManager"
	"shortener/internal/core/services/linkShortener"
	"shortener/internal/storage/postgres"
	"shortener/internal/storage/postgres/repository"
	"shortener/internal/transport/rest"
	"shortener/internal/transport/rest/handler/completeAdHandler"
	"shortener/internal/transport/rest/handler/getLinkClicksHandler"
	"shortener/internal/transport/rest/handler/getUserLinksHandler"
	"shortener/internal/transport/rest/handler/openLinkHandler"
	"shortener/internal/transport/rest/handler/shortLinkHandler"
	"shortener/internal/transport/rest/mw"
)

type FakePayer struct{}

func (f FakePayer) Pay(ctx context.Context, clickId int64) error {
	fmt.Println("oay")
	return nil
}

type FakeAdProvider struct{}

func (f FakeAdProvider) GetAdByMetadata(ctx context.Context, metadata model.ClickMetadata) (int64, error) {
	return 1, nil
}

func main() {
	fmt.Println("Hello World")

	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, cancel, err := postgres.NewPgxPool(cfg.PostgresURL)
	if err != nil {
		log.Error("unable to connect to postgres", mw.ErrAttr(err))
		cancel()
	}

	defer cancel()

	repo := repository.New(db)

	aliasGenerator := generator.New([]rune("i"), cfg.Service.GeneratedAliasLength) //todo alp
	shortener, err := linkShortener.New(repo, aliasGenerator, 3)                   //todo triesToGenerate
	linkService := linkManager.New(repo)

	adViewManager := adViewer.New(repo, &FakePayer{}, &FakeAdProvider{})
	onCompleteErrs := adViewManager.OnCompleteErrs()

	go func() {
		for err := range onCompleteErrs {
			log.Error(err.Error())
		}
	}()

	handlers := rest.Handlers{
		ShortLink:     shortLinkHandler.New(shortener),
		GetUserLinks:  getUserLinksHandler.New(linkService),
		GetLinkClicks: getLinkClicksHandler.New(linkService),
		OpenLink:      openLinkHandler.New(adViewManager),
		CompleteAd:    completeAdHandler.New(adViewManager),
	}

	server := rest.New(handlers, cfg.Auth, cfg.HTTPServer, log)

	if err := server.ListenAndServe(); err != nil {
		log.Error("Unable to start server: ", err)
		os.Exit(1) //todo
	}
}
