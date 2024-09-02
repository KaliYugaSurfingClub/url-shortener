package handler

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"link_shortener/internal/http/middlewares/mwLogger"
	"link_shortener/internal/lib/Api"
	"link_shortener/internal/lib/random"
	"link_shortener/internal/lib/sl"
	"link_shortener/internal/storage"
	"log/slog"
	"net/http"
)

type request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"omitempty,alphanumunicode"`
}

type response struct {
	Api.Response
	Alias string `json:"alias"`
}

type urlSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func Save(saver urlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mwLogger.GetCtxLog(r.Context(), "handlers.save.New")

		var req request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode request", sl.ErrorAttr(err))
			render.JSON(w, r, Api.Error("cannot decode request"))
			return
		}
		log = log.With(slog.Any("request", req))
		log.Debug("request body decoded")

		if errs := validator.New().Struct(req); errs != nil {
			log.Error("validation error", slog.String("error", errs.Error()))
			errForUser := Api.ValidationError(errs.(validator.ValidationErrors))
			render.JSON(w, r, errForUser)
			return
		}
		log.Debug("request is valid")

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(4)
		}
		log = log.With(slog.String("alias", alias))

		err := saver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrAliasExists) && req.Alias != "" {
			log.Info("alias was not added to database because it already exists")
			render.JSON(w, r, Api.Error(storage.ErrAliasExists.Error()))
			return
		}
		if err != nil {
			log.Error("failed to add url", slog.String("error", err.Error()))
			render.JSON(w, r, Api.Error("failed to add url"))
			return
		}

		log.Info("url added")

		render.JSON(w, r, response{
			Response: Api.Ok(),
			Alias:    req.Alias,
		})
	}
}

func saveWithoutAlias(w http.ResponseWriter, r *http.Request, saver urlSaver, log *slog.Logger, url string, alias string) {
	err := saver.SaveURL(url, alias)
	if errors.Is(err, storage.ErrAliasExists) {
		log.Info("alias was not added to database because it already exists")
		render.JSON(w, r, Api.Error(storage.ErrAliasExists.Error()))
		return
	}
	if err != nil {
		log.Error("failed to add url", slog.String("error", err.Error()))
		render.JSON(w, r, Api.Error("failed to add url"))
		return
	}
}

func saveWithAlias() {

}
