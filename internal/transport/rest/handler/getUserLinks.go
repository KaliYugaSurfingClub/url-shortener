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

type LinksProvider interface {
	GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error)
	GetCountByUserId(ctx context.Context, userId int64, params model.LinkFilter) (int64, error)
}

func GetUserLinks(provider LinksProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.GetUserLinks")

		//userId, _ := mw.ExtractUserID(r.Context())

		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode body", mw.ErrAttr(err))
			render.JSON(w, r, rest.NewErrorResponse(err))
			return
		}

		//var resp response
		//totalCount, err := provider.GetCountByUserId(userId, params)

	}
}

//func jsonToParams(req request) model.GetLinksParams {
//	//return model.GetLinksParams{
//	//	Filter: model.LinkFilter{
//	//		Type: model.LinkType(slices.Index(types, req.Type)),
//	//
//	//	}
//	//}
//}
