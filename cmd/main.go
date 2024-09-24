package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"url_shortener/core/generator"
	"url_shortener/core/managers/aliasManager"
	"url_shortener/core/managers/redirectManager"
	"url_shortener/core/model"
	"url_shortener/storage/sqlite/clickRepo"
	"url_shortener/storage/sqlite/linkRepo"
	"url_shortener/storage/transaction"
)

type FakeUserStore struct{}

func (u *FakeUserStore) AddToBalance(ctx context.Context, id int64, payment int) error {
	fmt.Println("balance increasing")
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

	aliasGenerator := generator.New([]rune("1"), 1)
	s, _ := aliasManager.New(linkStore, aliasGenerator, 10)

	if _, err := s.Save(context.Background(), &model.Link{Original: "orig"}); err != nil {
		log.Fatal(err)
	}

	r := redirectManager.New(linkStore, clickStore, &FakeUserStore{}, transactor)

	//create click and mark that AD was started
	original, clickId, userId, err := r.Start(context.Background(), "1", &model.ClickMetadata{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original: ", original)

	//reward creator of link and mark that AD was watched
	err = r.End(context.Background(), clickId, userId)
	if err != nil {
		log.Fatal(err)
	}
}
