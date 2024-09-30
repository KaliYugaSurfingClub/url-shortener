package handler

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest"
	"shortener/internal/transport/rest/mw"
	"time"
)

type LinksProvider interface {
	GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error)
	GetCountByUserId(ctx context.Context, userId int64, params model.LinkFilter) (int64, error)
}

func GetUserLinks(provider LinksProvider) http.HandlerFunc {
	type request struct {
		Type   string `json:"type"`
		SortBy string `json:"sortBy"`
		Order  string `json:"order"`
		Page   int    `json:"page"`
		Size   int    `json:"size"`
	}

	type link struct {
		Id             int64      `json:"id"`
		CreatedBy      int64      `json:"createdBy"`
		Original       string     `json:"original"`
		Alias          string     `json:"alias"`
		CustomName     string     `json:"customName"`
		ClicksCount    int64      `json:"clicksCount"`
		LastAccessTime *time.Time `json:"lastAccessTime,omitempty"`
		ExpirationDate *time.Time `json:"expirationDate,omitempty"`
		ClicksToExpire *int64     `json:"clicksToExpire,omitempty"`
		Archived       bool       `json:"archived"`
		CreatedAt      time.Time  `json:"createdAt"`
	}

	type response struct {
		TotalCount int64  `json:"totalCount"`
		Links      []link `json:"links,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

		userId, _ := mw.ExtractUserID(r.Context())

		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode body", mw.ErrAttr(err))
			render.JSON(w, r, rest.NewErrorResponse(err))
			return
		}

		params := model.GetLinksParams{Pagination: model.Pagination{Size: 10, Page: 1}}

		var resp response

		totalCount, err := provider.GetCountByUserId(r.Context(), userId, params.Filter)
		if err != nil {
			log.Error("cannot get count of user links", mw.ErrAttr(err))
			render.JSON(w, r, rest.NewErrorResponse(err))
			return
		}

		resp.TotalCount = totalCount
		resp.Links = make([]link, 0, resp.TotalCount)

		links, err := provider.GetByUserId(r.Context(), userId, params)
		if err != nil {
			log.Error("cannot get user links", mw.ErrAttr(err))
			render.JSON(w, r, rest.NewErrorResponse(err))
			return
		}

		//todo add go generate mapper
		for _, l := range links {
			resp.Links = append(resp.Links, link{
				Id:             l.Id,
				CreatedBy:      l.CreatedBy,
				Original:       l.Original,
				Alias:          l.Alias,
				CustomName:     l.CustomName,
				ClicksCount:    l.ClicksCount,
				LastAccessTime: l.LastAccessTime,
				ExpirationDate: l.ExpirationDate,
				ClicksToExpire: l.ClicksToExpire,
				Archived:       l.Archived,
				CreatedAt:      l.CreatedAt,
			})
		}

		render.JSON(w, r, resp)
	}
}
