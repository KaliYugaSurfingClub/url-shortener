package postgres

//
//import (
//	"context"
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"log"
//	"os"
//	"shortener/internal/core/model"
//	"shortener/internal/storage/postgres/linkRepo"
//	"testing"
//	"time"
//)
//
//var (
//	dbPath    = "./test.db"
//	logsPath  = "allSelectsWithParams.txt"
//	repo      *linkRepo.LinkRepo
//	creatorId int64 = 1
//)
//
//func init() {
//	db, err := sqlx.Open("sqlite3", dbPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	repo = linkRepo.New(db)
//}
//
//func GetAllTypeLinks() map[model.LinkType]string {
//	return map[model.LinkType]string{
//		model.TypeAny:      "TypeAny",
//		model.TypeActive:   "TypeActive",
//		model.TypeExpired:  "TypeExpired",
//		model.TypeArchived: "TypeArchived",
//	}
//}
//
//func GetAllConstraintLinks() map[model.LinkConstraints]string {
//	return map[model.LinkConstraints]string{
//		model.ConstraintAny:    "ConstraintAny",
//		model.ConstraintClicks: "ConstraintClicks",
//		model.ConstraintDate:   "ConstraintDate",
//		model.ConstraintWith:   "ConstraintWith",
//	}
//}
//
//func GetAllColumnLinks() map[model.SortByLink]string {
//	return map[model.SortByLink]string{
//		model.SortByCreatedAt:       "ColumnCreatedAt",
//		model.SortByCustomName:      "ColumnCustomName",
//		model.SortByClicksCount:     "ColumnClicksCount",
//		model.SortByLastAccess:      "ColumnLastAccess",
//		model.SortByExpirationDate:  "ColumnTimeToExpire",
//		model.SortByLeftClicksCount: "ColumnClicksToExpire",
//	}
//}
//
//func GetAllOrders() map[model.Order]string {
//	return map[model.Order]string{
//		model.Desc: "Desc",
//		model.Asc:  "Asc",
//	}
//}
//
//type testCase struct {
//	params model.GetLinksParams
//	name   string
//}
//
//func generateCombinations() []testCase {
//	var combinations []testCase
//
//	for t, ts := range GetAllTypeLinks() {
//		for c, cs := range GetAllConstraintLinks() {
//			for col, cols := range GetAllColumnLinks() {
//				for or, ors := range GetAllOrders() {
//					combinations = append(combinations, testCase{
//						name: ts + " " + cs + " " + cols + " " + ors,
//						params: model.GetLinksParams{
//							Type:        t,
//							Constraints: c,
//							SortBy:      col,
//							Order:       or,
//						},
//					})
//				}
//			}
//		}
//	}
//
//	return combinations
//}
//
//func printSortedFilteredLinks(name string, links []*model.Link, w io.Writer) error {
//	if _, err := fmt.Fprintf(w, "%s - %d\n", name, len(links)); err != nil {
//		return err
//	}
//
//	for _, link := range links {
//		var leftClicks int64
//
//		if link.ClicksToExpiration != nil {
//			leftClicks = *link.ClicksToExpiration - link.ClicksCount
//		}
//
//		if link.ExpirationDate == nil {
//			link.ExpirationDate = &time.Time{}
//		}
//
//		if link.ClicksToExpiration == nil {
//			var z int64
//			link.ClicksToExpiration = &z
//		}
//
//		if link.LastAccessTime == nil {
//			link.LastAccessTime = &time.Time{}
//		}
//
//		_, err := fmt.Fprintf(w,
//			"{\n\tCustomName - %s\n"+
//				"\tClicksCount - %d\n"+
//				"\tClicksToExpiration - %d\n"+
//				"\tLeftClicksCount - %d\n"+
//				"\tLastAccessTime %v\n"+
//				"\tExpirationDate %v\n"+
//				"\tArchived %t\n"+
//				"\tCreatedAt %v\n"+
//				"}\n",
//			link.CustomName,
//			link.ClicksCount,
//			*link.ClicksToExpiration,
//			leftClicks,
//			*link.LastAccessTime,
//			*link.ExpirationDate,
//			link.Archived,
//			link.CreatedAt,
//		)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func TestSelectAll(t *testing.T) {
//	t.Logf("VISIT %q TO CHECK RESULTS OF EACH QUERY THIS TEST", logsPath)
//
//	file, _ := os.OpenFile(logsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
//
//	for _, test := range generateCombinations() {
//		t.Run(test.name, func(t *testing.T) {
//			count, err := repo.GetCount(context.Background(), creatorId, test.params)
//			require.NoError(t, err)
//
//			links, err := repo.GetByUserId(context.Background(), creatorId, test.params)
//			require.NoError(t, err)
//
//			if err = printSortedFilteredLinks(test.name, links, file); err != nil {
//				log.Fatal(err)
//			}
//
//			assert.Equal(t, count, int64(len(links)), "count method and len of got links")
//		})
//	}
//}
