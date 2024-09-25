package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"shortener/internal/storage/postgres/linkRepo"
	"time"
)

type FakeUserStore struct{}

func (u *FakeUserStore) AddToBalance(ctx context.Context, id int64, payment int) error {
	fmt.Println("balance increasing")
	return nil
}

func main() {
	//todo literal
	dbURL := "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"

	poolCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL: ", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatal("Unable to create connection pool: ", err)
	}

	defer db.Close()

	//clickStore := clickRepo.New(db)
	//
	//id, err := clickStore.Save(context.Background(), &model.Click{LinkId: 1})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(id)
	//
	//err = clickStore.UpdateStatus(context.Background(), id, model.AdClosed)
	//if err != nil {
	//	log.Fatal(err)
	//}

	linkStore := linkRepo.New(db)

	link, err := linkStore.GetActiveByAlias(context.Background(), "abcd")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(link)

	//linkStore := linkRepo.New(db)
	//transactor := transaction.NewTransactor(db)

	//aliasGenerator := generator.New([]rune("1"), 1)
	//s, _ := aliasManager.New(linkStore, aliasGenerator, 10)
	//
	//if _, err := s.Save(context.Background(), &model.Link{Original: "orig"}); err != nil {
	//	log.Fatal(err)
	//}
	//
	//r := redirectManager.New(linkStore, clickStore, &FakeUserStore{}, transactor)
	//
	////create click and mark that AD was started
	//original, clickId, userId, err := r.Start(context.Background(), "1", &model.ClickMetadata{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("Original: ", original)
	//
	////reward creator of link and mark that AD was watched
	//err = r.End(context.Background(), clickId, userId)
	//if err != nil {
	//	log.Fatal(err)
	//}

	time.Sleep(3 * time.Second)
}
