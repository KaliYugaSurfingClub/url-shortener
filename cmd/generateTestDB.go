package main

import (
	"context"
	"errors"
	"github.com/jaswdr/faker"
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
	"time"
	"url_shortener/core"
	"url_shortener/core/model"
	"url_shortener/storage/sqlite/linkRepo"
)

const CreatorId = 1

func oneFromTwo[T any](a T, b T) T {
	if rand.Float32() < 0.5 {
		return a
	}

	return b
}

func generateLink(fake *faker.Faker) model.Link {
	maxClicks := oneFromTwo(
		int64(rand.Int())%1024,
		model.UnlimitedClicks,
	)

	return model.Link{
		CreatedBy: CreatorId,
		Original:  fake.Lorem().Word(),
		Alias:     fake.Lorem().Word(),
		ExpirationDate: oneFromTwo(
			fake.Time().TimeBetween(time.Now(), time.Now().Add(10000*time.Hour)),
			model.NoExpireDate,
		),
		ClicksCount: oneFromTwo(
			rand.Int63(),
			maxClicks+rand.Int63(),
		),
		LastAccessTime: oneFromTwo(
			fake.Time().TimeBetween(time.Now(), time.Now().Add(10000*time.Hour)),
			time.Now(),
		),
		MaxClicks: maxClicks,
	}
}

func main() {
	db, err := sqlx.Open("sqlite3", "C:\\Users\\leono\\Desktop\\prog\\go\\url_shortener\\test.db")
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(323213)
	fake := faker.New()

	repo := linkRepo.New(db)

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)

		_, err := repo.Save(context.Background(), generateLink(&fake))

		if err != nil && !errors.Is(err, core.ErrAliasExists) {
			log.Fatal(err, " saving error")
		}
	}
}
