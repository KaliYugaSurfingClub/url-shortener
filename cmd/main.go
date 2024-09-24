package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
	"url_shortener/core/generator"
	"url_shortener/core/managers/aliasManager"
	"url_shortener/core/managers/redirectManager"
	"url_shortener/core/model"
	"url_shortener/storage/sqlite/clickRepo"
	"url_shortener/storage/sqlite/linkRepo"
	"url_shortener/storage/transaction"
)

func main() {
	db, err := sqlx.Open("sqlite3", "C:\\Users\\leono\\Desktop\\prog\\go\\url_shortener\\storage.db")
	if err != nil {
		log.Fatal(err)
	}

	clickStore := clickRepo.New(db)
	linkStore := linkRepo.New(db)
	transactor := transaction.NewTransactor(db)

	aliasGenerator := generator.New([]rune("e"), 1)
	s, _ := aliasManager.New(linkStore, aliasGenerator, 10)

	if _, err := s.Save(context.Background(), &model.Link{Original: "orig"}); err != nil {
		log.Fatal(err)
	}

	r := redirectManager.New(linkStore, clickStore, transactor)

	err = r.HandleClick(context.Background(), "e", &model.ClickMetadata{})
	if err != nil {
		log.Fatal(err)
	}

	err = r.HandleClick(context.Background(), "e", &model.ClickMetadata{})
	if err != nil {
		log.Fatal(err)
	}

	//todo sleep for test
	time.Sleep(10 * time.Second)
}
