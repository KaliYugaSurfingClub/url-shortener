package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
	"url_shortener/core/model"
	"url_shortener/storage/sqlite"
)

func main() {
	db, err := sqlx.Open("sqlite3", "storage.db")
	if err != nil {
		log.Fatal(err)
	}

	clickStore := sqlite.NewClickRepo(db)
	linkStore := sqlite.NewLinkRepo(db)
	//transactor := transaction.NewTransactor(db)

	_, err = linkStore.Save(context.Background(), model.Link{
		CreatedBy: 1, Original: "abc", Alias: "123", LastAccess: time.Now(), ExpireDate: model.NoExpireDate,
		MaxClicks: model.UnlimitedClicks,
	})

	if err != nil {
		log.Fatal(err)
	}

	//r := redirect.New(linkStore, linkStore, clickStore, func(string) {}, transactor)
	//if err = r.To(context.Background(), "123"); err != nil {
	//	log.Fatal(err)
	//}
	//
	//aliasGenerator := generator.New([]rune("ab"), 1)
	//
	//s := alias.New(linkStore, aliasGenerator, 10)
	//if al, err := s.Save(context.Background(), "123", ""); err != nil {
	//	log.Fatal(err)
	//} else {
	//	fmt.Printf(al)
	//}

}
