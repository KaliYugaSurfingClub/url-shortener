package sqlite

import (
	"github.com/jaswdr/faker"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"testing"
	"time"
	"url_shortener/core/model"
	"url_shortener/storage/entity"
)

func oneFromTwo[T any](a T, b T) T {
	if rand.Float32() < 0.5 {
		return a
	}

	return b
}

func generateLinks(fake *faker.Faker) model.Link {
	clicksToExpiration := oneFromTwo(
		int64(rand.Int())%1024,
		model.UnlimitedClicks,
	)

	var clicksCount int64
	if clicksToExpiration != model.UnlimitedClicks {
		clicksCount = int64(rand.Intn(int(clicksToExpiration) + 100))
	}

	return model.Link{
		CreatedBy:          creatorId,
		Original:           fake.Lorem().Word(),
		Alias:              fake.Lorem().Word(),
		CustomName:         fake.Lorem().Word(),
		ClicksCount:        clicksCount,
		LastAccessTime:     fake.Time().TimeBetween(time.Now().Add(-1000*time.Hour), time.Now().Add(1000*time.Hour)),
		ClicksToExpiration: clicksToExpiration,
		Archived:           oneFromTwo(true, false),

		ExpirationDate: oneFromTwo(
			fake.Time().TimeBetween(time.Now().Add(-1000*time.Hour), time.Now().Add(1000*time.Hour)),
			model.NoExpireDate,
		),
	}
}

func TestMain(m *testing.M) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec(`DELETE FROM link`)

	fake := faker.New()
	for i := 0; i < 50; i++ {
		time.Sleep(1 * time.Second)

		link := generateLinks(&fake)

		//todo
		db.Exec(`
			INSERT INTO link(
				created_by, original, alias, custom_name, clicks_count, 
				last_access_time, expiration_date, clicks_to_expiration, archived) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?
			)`,
			entity.CreatedByToSql(link.CreatedBy),
			link.Original,
			link.Alias,
			link.CustomName,
			link.ClicksCount,
			link.LastAccessTime,
			entity.ExpirationDateToSql(link.ExpirationDate),
			entity.ClicksToExpirationToSql(link.ClicksToExpiration),
			link.Archived,
		)
	}
}
