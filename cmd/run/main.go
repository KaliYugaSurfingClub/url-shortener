package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"shortener/internal/core/generator"
	"shortener/internal/core/managers/aliasManager"
	"shortener/internal/core/managers/redirectManager"
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres/clickRepo"
	"shortener/internal/storage/postgres/linkRepo"
	"shortener/internal/storage/transaction"
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
		fmt.Println("Unable to parse DATABASE_URL: ", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		fmt.Println("Unable to create connection pool: ", err)
	}

	defer db.Close()

	clickStore := clickRepo.New(db)

	//id, err := clickStore.Save(context.Background(), &model.Click{LinkId: 1})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(id)
	//
	//err = clickStore.UpdateStatus(context.Background(), id, model.AdClosed)
	//if err != nil {
	//	fmt.Println(err)
	//}

	linkStore := linkRepo.New(db)

	//link := &model.Link{
	//	CreatedBy: 1,
	//	Original:  "abcr",
	//	Alias:     "dsaads",
	//}
	//
	//id, err = linkStore.Save(context.Background(), link)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(id)
	//
	//link, err = linkStore.GetActiveByAlias(context.Background(), "abcd")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(link)
	//
	//links, err := linkStore.GetByUserId(context.Background(), 1, model.GetLinksParams{})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(links)
	//
	//count, err := linkStore.GetCount(context.Background(), 1, model.LinkFilter{})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(count)
	//
	//err = linkStore.UpdateLastAccess(context.Background(), id, time.Now())
	//if err != nil {
	//	fmt.Println(err)
	//}

	transactor := transaction.NewTransactor(db)

	aliasGenerator := generator.New([]rune("1"), 1)
	s, _ := aliasManager.New(linkStore, aliasGenerator, 10)

	if _, err := s.Save(context.Background(), &model.Link{Original: "orig", CreatedBy: 1}); err != nil {
		fmt.Println(err)
	}

	r := redirectManager.New(linkStore, clickStore, &FakeUserStore{}, transactor)

	//create click and mark that AD was started
	link, clickId, err := r.Start(context.Background(), "1", &model.ClickMetadata{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Original: ", link.Original)

	//reward creator of link and mark that AD was watched
	err = r.End(context.Background(), clickId, link.CreatedBy)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(3 * time.Second)
}
