package main

import (
	"context"
	"fmt"
	"github.com/KaliYugaSurfingClub/pkg/mw"
	"link-service/internal/config"
	"link-service/internal/core/generator"
	"link-service/internal/core/model"
	"link-service/internal/core/services/adViewer"
	"link-service/internal/core/services/linkManager"
	"link-service/internal/core/services/linkShortener"
	"link-service/internal/storage/postgres"
	"link-service/internal/storage/postgres/repository"
	"link-service/internal/transport/rest/handler/completeAdHandler"
	"link-service/internal/transport/rest/handler/getLinkClicksHandler"
	"link-service/internal/transport/rest/handler/getUserLinksHandler"
	"link-service/internal/transport/rest/handler/openLinkHandler"
	"link-service/internal/transport/rest/handler/shortLinkHandler"
	"link-service/internal/transport/rest/server"
	"log/slog"
	"os"
)

type FakePayer struct{}

func (f FakePayer) Pay(ctx context.Context, clickId int64) error {
	fmt.Println("pay")
	return nil
}

type FakeAdProvider struct{}

func (f FakeAdProvider) GetAdByMetadata(ctx context.Context, metadata model.ClickMetadata) (int64, error) {
	return 1, nil
}

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("logger enabled")

	db, cancel, err := postgres.NewPgxPool(cfg.PostgresURL)
	if err != nil {
		log.Error("unable to connect to postgres", mw.ErrAttr(err))
		cancel()
		os.Exit(1)
	}

	log.Info("connection with postgres established")

	defer cancel()

	repo := repository.New(db)

	aliasGenerator := generator.New([]rune(cfg.Service.Alp), cfg.Service.GeneratedAliasLength)
	shortener := linkShortener.New(repo, aliasGenerator, cfg.Service.TriesToGenerate)
	linkService := linkManager.New(repo)

	adViewManager := adViewer.New(repo, &FakePayer{}, &FakeAdProvider{})
	onCompleteErrs := adViewManager.OnCompleteErrs()

	go func() {
		for err := range onCompleteErrs {
			log.Error(err.Error())
		}
	}()

	handlers := server.Handlers{
		ShortLink:     shortLinkHandler.New(shortener),
		GetUserLinks:  getUserLinksHandler.New(linkService),
		GetLinkClicks: getLinkClicksHandler.New(linkService),
		OpenLink:      openLinkHandler.New(adViewManager),
		CompleteAd:    completeAdHandler.New(adViewManager),
	}

	log.Info(fmt.Sprintf("server starting on %s", cfg.HTTPServer.IP+":"+cfg.HTTPServer.Port))
	server := server.New(handlers, cfg.Auth, cfg.HTTPServer, log)
	if err := server.ListenAndServe(); err != nil {
		log.Error("Unable to start server: ", err)
		cancel()
		os.Exit(1)
	}
}
