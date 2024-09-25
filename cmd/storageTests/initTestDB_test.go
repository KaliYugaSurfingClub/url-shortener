package postgres

//import (
//	"github.com/jaswdr/faker"
//	"github.com/jmoiron/sqlx"
//	"log"
//	"math/rand"
//	"shortener/internal/core/model"
//	"shortener/internal/storage/postgres/entity"
//	"testing"
//	"time"
//)
//
//func oneFromTwo[T any](a T, b T) T {
//	if rand.Float32() < 0.5 {
//		return a
//	}
//
//	return b
//}
//
//func generateLinks(fake *faker.Faker) model.Link {
//	res := model.Link{
//		CreatedBy:  &creatorId,
//		Original:   fake.Lorem().Word(),
//		Alias:      fake.Lorem().Word(),
//		CustomName: fake.Lorem().Word(),
//		Archived:   oneFromTwo(true, false),
//	}
//
//	if rand.Int()%2 == 0 {
//		clicksToExpiration := int64(rand.Int()) % 1024
//		res.ClicksToExpiration = &clicksToExpiration
//	}
//
//	if res.ClicksToExpiration != nil {
//		res.ClicksCount = int64(rand.Intn(int(*res.ClicksToExpiration) + 100))
//	}
//
//	return res
//}
//
//func TestDBInit(t *testing.T) {
//	db, err := sqlx.Open("sqlite3", dbPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	db.Exec(`DELETE FROM link`)
//
//	fake := faker.New()
//	for i := 0; i < 50; i++ {
//		time.Sleep(1 * time.Second)
//
//		link := generateLinks(&fake)
//
//		query := `
//			INSERT INTO link(created_by, original, alias, custom_name, clicks_count,
//							 last_access_time, expiration_date, clicks_to_expiration, archived)
//			VALUES (:created_by, :original, :alias, :custom_name, :clicks_count,
//					:last_access_time, :expiration_date, :clicks_to_expiration, :archived)
//		`
//
//		db.NamedExec(query, entity.ModelToLink(&link))
//	}
//}
