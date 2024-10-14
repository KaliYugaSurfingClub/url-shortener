package openShortenedHandler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net"
	"net/http"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/transport/rest/mw"
	"strconv"
	"strings"
	"time"
)

type AdPageProvider interface {
	GetAdPage(ctx context.Context, alias string, metadata model.ClickMetadata) (*model.AdPage, error)
}

type Handler struct {
	adPageProvider   AdPageProvider
	adTypeToTemplate map[model.AdType]*template.Template
	AdSourceHost     string
	CallBackURL      string
}

func New(
	adPageProvider AdPageProvider,
	adTypeToTemplate map[model.AdType]*template.Template,
	AdSourceHost string,
	CallBackURL string,
) *Handler {
	return &Handler{
		adTypeToTemplate: adTypeToTemplate,
		adPageProvider:   adPageProvider,
		AdSourceHost:     AdSourceHost,
		CallBackURL:      CallBackURL,
	}
}

type pageVars struct {
	Original    string
	CallBackURL string
	AdURL       string
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

	alias := chi.URLParam(r, "alias")

	metadata := model.ClickMetadata{
		UserAgent:  r.UserAgent(),
		AccessTime: time.Now(),
		IP:         getClientIP(r),
	}

	adPage, err := h.adPageProvider.GetAdPage(r.Context(), alias, metadata)
	if errors.Is(err, core.ErrLinkNotFound) {
		log.Info(err.Error())
		w.Write([]byte("not found")) //todo
		return
	}
	if err != nil {
		log.Info(err.Error())
		w.Write([]byte("internal error")) //todo
		return
	}

	pv := pageVars{
		Original:    adPage.Original,
		AdURL:       h.AdSourceHost + "/" + strconv.FormatInt(adPage.AdSourceId, 10),
		CallBackURL: h.CallBackURL + "/" + strconv.FormatInt(adPage.ClickId, 10),
	}

	err = h.adTypeToTemplate[adPage.AdType].Execute(w, pv)
	if err != nil {
		log.Info(err.Error())
		w.Write([]byte("internal error")) //todo
		return
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
