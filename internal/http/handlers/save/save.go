package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"link_shortener/internal/lib/random"
	"link_shortener/internal/lib/response"
	"link_shortener/internal/storage"
	"log/slog"
	"net/http"
)

//todo maybe refactor directories

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"omitempty,alphanumunicode"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		var req Request

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode request", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("cannot decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		if errs := validator.New().Struct(req); errs != nil {
			log.Error("validation error", slog.String("error", errs.Error()))
			render.JSON(w, r, response.ValidationError(errs.(validator.ValidationErrors)))
			return
		}

		//todo random alias may already be in database
		if req.Alias == "" {
			//todo literal should not be here.
			req.Alias = random.NewRandomString(4)
		}

		log = log.With(
			slog.String("alias", req.Alias),
			slog.String("url", req.URL),
		)

		err := saver.SaveURL(req.URL, req.Alias)
		if errors.Is(err, storage.ErrAliasExists) {
			log.Info("alias was not added to database because it already exists")
			render.JSON(w, r, response.Error(storage.ErrAliasExists.Error()))
			return
		}
		if err != nil {
			log.Error("failed to add url", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		log.Info("url added")

		render.JSON(w, r, Response{
			Response: response.Ok(),
			Alias:    req.Alias,
		})
	}
}
