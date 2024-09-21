package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"url_shortener/core/alias"
	"url_shortener/core/redirect"
	"url_shortener/generator"
	"url_shortener/storage"
	"url_shortener/storage/transaction"
)

type fakeUpdater struct{}

func (fs *fakeUpdater) UpdateLastAccess(ctx context.Context, id int, t time.Time) error {
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = storage.InitTables(db); err != nil {
		log.Fatal(err)
	}

	transactor := transaction.NewTransactor(db)
	linkStore := storage.NewLinkRepo(db)
	clickStore := storage.NewClickRepo(db)

	r := redirect.New(linkStore, linkStore, clickStore, func(string) {}, transactor)
	if err = r.To(context.Background(), "123"); err != nil {
		log.Fatal(err)
	}

	aliasGenerator := generator.New([]rune("ab"), 1)

	s := alias.New(linkStore, aliasGenerator, 10)
	if al, err := s.Save(context.Background(), "123", ""); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf(al)
	}

}
