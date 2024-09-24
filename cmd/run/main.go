package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"shortener/internal/core/generator"
	"shortener/internal/core/managers/aliasManager"
	"shortener/internal/core/managers/redirectManager"
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres/clickRepo"
	"shortener/internal/storage/postgres/linkRepo"
	"shortener/internal/storage/transaction"
)

type FakeUserStore struct{}

func (u *FakeUserStore) AddToBalance(ctx context.Context, id int64, payment int) error {
	fmt.Println("balance increasing")
	return nil
}

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	defer db.Close()

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
