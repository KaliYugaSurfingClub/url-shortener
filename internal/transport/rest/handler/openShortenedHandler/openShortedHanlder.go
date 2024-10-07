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
	"strings"
	"time"
)

type recorder interface {
	OnClick(ctx context.Context, alias string, metadata model.ClickMetadata) (string, int64, error)
}

type adProvider interface {
	Get(ctx context.Context) (string, error)
}

type Handler struct {
	recorder       recorder
	videoURL       string
	adPageTemplate *template.Template
}

func New(
	recorder recorder,
	videoURL string,
	adPageTemplate *template.Template,
) *Handler {
	return &Handler{
		recorder:       recorder,
		videoURL:       videoURL,
		adPageTemplate: adPageTemplate,
	}
}

type PageVariables struct {
	ClickId  int64
	Original string
	VideoURL string
}

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) {
	log := mw.ExtractLog(r.Context(), "transport.rest.openShortenedHandler")

	alias := chi.URLParam(r, "alias")
	metadata := model.ClickMetadata{
		UserAgent:  r.UserAgent(),
		AccessTime: time.Now(),
		IP:         getClientIP(r),
	}

	original, clickId, err := h.recorder.OnClick(r.Context(), alias, metadata)
	if errors.Is(err, core.ErrLinkNotFound) {
		log.Info(err.Error())
		w.Write([]byte("not found")) //todo
		return
	}

	variables := PageVariables{
		Original: original,
		ClickId:  clickId,
		VideoURL: h.videoURL,
	}

	err = h.adPageTemplate.Execute(w, variables)
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
