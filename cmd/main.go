package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"url_shortener/core/Managers/aliasManager"
	"url_shortener/core/Managers/redirectManager"
	"url_shortener/core/model"
	"url_shortener/generator"
	"url_shortener/storage/sqlite/clickRepo"
	"url_shortener/storage/sqlite/linkRepo"
	"url_shortener/storage/transaction"
)

type fakeReward struct{}

func (t *fakeReward) TransferReward(userId int64) error {
	return nil
}

func main() {
	db, err := sqlx.Open("sqlite3", "C:\\Users\\leono\\Desktop\\prog\\go\\url_shortener\\storage.db")
	if err != nil {
		log.Fatal(err)
	}

	clickStore := clickRepo.New(db)
	linkStore := linkRepo.New(db)
	transactor := transaction.NewTransactor(db)

	r := redirectManager.New(linkStore, linkStore, clickStore, &fakeReward{}, transactor)

	original, err := r.Process(context.Background(), "123", model.Click{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(original)

	aliasGenerator := generator.New([]rune("ab"), 1)

	s, _ := aliasManager.New(linkStore, aliasGenerator, 10)
	if al, err := s.Save(context.Background(), model.Link{}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf(al)
	}

}
