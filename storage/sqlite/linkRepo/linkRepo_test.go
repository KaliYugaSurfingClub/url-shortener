package linkRepo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
	"log"
	"strconv"
	"testing"
	"time"
	"url_shortener/core/model"
)

func GetAllTypeLinks() map[model.TypeLink]string {
	return map[model.TypeLink]string{
		model.TypeAny:     "TypeAny",
		model.TypeActual:  "TypeActual",
		model.TypeExpired: "TypeExpired",
	}
}

func GetAllConstraintLinks() map[model.ConstraintLink]string {
	return map[model.ConstraintLink]string{
		model.ConstraintAny:            "ConstraintAny",
		model.ConstraintMaxClicks:      "ConstraintMaxClicks",
		model.ConstraintExpirationDate: "ConstraintExpirationDate",
		model.ConstraintWith:           "ConstraintWith",
	}
}

func GetAllColumnLinks() map[model.ColumnLink]string {
	return map[model.ColumnLink]string{
		model.ColumnCreatedAt:      "ColumnCreatedAt",
		model.ColumnAlias:          "ColumnAlias",
		model.ColumnClicksCount:    "ColumnClicksCount",
		model.ColumnLastAccess:     "ColumnLastAccess",
		model.ColumnTimeToExpire:   "ColumnTimeToExpire",
		model.ColumnClicksToExpire: "ColumnClicksToExpire",
	}
}

func GetAllOrders() map[model.Order]string {
	return map[model.Order]string{
		model.Desc: "Desc",
		model.Asc:  "Asc",
	}
}

type testCase struct {
	params model.GetLinksParams
	name   string
}

func generateCombinations() []testCase {
	var combinations []testCase

	for t, ts := range GetAllTypeLinks() {
		for c, cs := range GetAllConstraintLinks() {
			for col, cols := range GetAllColumnLinks() {
				for o, os := range GetAllOrders() {
					combinations = append(combinations, testCase{
						name: ts + " " + cs + " " + cols + " " + os,
						params: model.GetLinksParams{
							Type:        t,
							Constraints: c,
							Column:      col,
							Order:       o,
						},
					})
				}
			}
		}
	}

	return combinations
}

func timesAreOrdered(a, b time.Time, order model.Order) bool {
	if order == model.Asc {
		return a.Before(b) || a.Equal(b)
	}

	return a.After(b) || a.Equal(b)
}

func valuesAreOrdered[T constraints.Ordered](a, b T, order model.Order) bool {
	if order == model.Asc {
		return a <= b
	}

	return a >= b
}

func SortValidate(links []*model.Link, params model.GetLinksParams) (bool, string, string) {
	for i := 1; i < len(links); i++ {
		switch params.Column {
		case model.ColumnCreatedAt:
			if !timesAreOrdered(links[i-1].CreatedAt, links[i].CreatedAt, params.Order) {
				return false, links[i-1].CreatedAt.String(), links[i].CreatedAt.String()
			}
		case model.ColumnClicksCount:
			if !valuesAreOrdered(links[i-1].ClicksCount, links[i].ClicksCount, params.Order) {
				return false,
					strconv.FormatInt(links[i-1].ClicksCount, 10),
					strconv.FormatInt(links[i].ClicksCount, 10)
			}
		case model.ColumnAlias:
			if !valuesAreOrdered(links[i-1].Alias, links[i].Alias, params.Order) {
				return false, links[i-1].Alias, links[i].Alias
			}
		case model.ColumnLastAccess:
			if !timesAreOrdered(links[i-1].LastAccessTime, links[i].LastAccessTime, params.Order) {
				return false, links[i-1].LastAccessTime.String(), links[i].LastAccessTime.String()
			}
		default:
			//todo add timeToExpire and clicksToExpire
		}
	}

	return true, "", ""
}

func filterByTypeValidate(links []*model.Link, params model.GetLinksParams) []model.Link {
	extraLinks := make([]model.Link, 0)

	for _, link := range links {
		switch params.Type {
		case model.TypeActual:
			if link.IsExpired() {
				extraLinks = append(extraLinks, *link)
			}
		case model.TypeExpired:
			if !link.IsExpired() {
				extraLinks = append(extraLinks, *link)
			}
		case model.TypeAny: //any one is suitable
		}
	}

	return extraLinks
}

func filterByConstraintsValidate(links []*model.Link, params model.GetLinksParams) []model.Link {
	extraLinks := make([]model.Link, 0)

	for _, link := range links {
		switch params.Constraints {
		case model.ConstraintWithout:
			if link.MaxClicks != model.UnlimitedClicks || link.ExpirationDate != model.NoExpireDate {
				extraLinks = append(extraLinks, *link)
			}
		case model.ConstraintWith:
			if link.MaxClicks == model.UnlimitedClicks || link.ExpirationDate == model.NoExpireDate {
				extraLinks = append(extraLinks, *link)
			}
		case model.ConstraintMaxClicks:
			if link.MaxClicks == model.UnlimitedClicks {
				extraLinks = append(extraLinks, *link)
			}
		case model.ConstraintExpirationDate:
			if link.ExpirationDate == model.NoExpireDate {
				extraLinks = append(extraLinks, *link)
			}
		case model.ConstraintAny: //any one is suitable
		}
	}

	return extraLinks
}

func TestNewClickRepoBasicCase(t *testing.T) {
	db, err := sqlx.Open("sqlite3", "C:\\Users\\leono\\Desktop\\prog\\go\\url_shortener\\test.db")
	if err != nil {
		log.Fatal(err)
	}

	repo := New(db)

	for _, test := range generateCombinations() {
		t.Run(test.name, func(t *testing.T) {
			links, err := repo.GetByUserId(context.Background(), 1, test.params)

			require.NoError(t, err)

			extraLinksByType := filterByTypeValidate(links, test.params)
			assert.Equal(
				t, 0, len(extraLinksByType),
				fmt.Sprintf("extraLinksByType %+v", extraLinksByType),
			)

			extraLinksByConstraints := filterByConstraintsValidate(links, test.params)
			assert.Equal(
				t, 0, len(extraLinksByConstraints),
				fmt.Sprintf("extraLinksByConstraints %+v", extraLinksByConstraints),
			)

			ok, a, b := SortValidate(links, test.params)
			assert.Equal(t, true, ok, fmt.Sprintf("%s %s are not ordered", a, b))
		})
	}
}
