package openLinkHandler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest"
	"shortener/internal/transport/rest/mw"
	"strings"
	"time"
)

type AdPageProvider interface {
	GetAdPage(ctx context.Context, alias string, metadata model.ClickMetadata) (*model.AdPage, error)
}

type response struct {
	AdType     string `json:"type"`
	ClickId    int64  `json:"click_id"`
	AdSourceId int64  `json:"ad_source_id"`
}

func New(adPageProvider AdPageProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")
		alias := chi.URLParam(r, "alias")

		metadata := model.ClickMetadata{
			UserAgent:  r.UserAgent(),
			AccessTime: time.Now(),
			IP:         getClientIP(r),
		}

		adPage, err := adPageProvider.GetAdPage(r.Context(), alias, metadata)
		if err != nil {
			rest.Error(w, log, err)
			return
		}

		rest.Ok(w, response{
			ClickId:    adPage.ClickId,
			AdSourceId: adPage.AdSourceId,
			AdType:     string(adPage.AdType),
		})
	}
}

func getClientIP(r *http.Request) net.IP {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		ip := strings.TrimSpace(ips[0])

		return net.ParseIP(ip)
	}

	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip, _, _ = net.SplitHostPort(ip)
	}

	return net.ParseIP(ip)
}
