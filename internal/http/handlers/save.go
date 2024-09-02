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

		if req.Alias != "" {
			saveWithoutAlias(w, r, saver, log, req.URL, req.Alias)
		} else {
			saveWithAlias(w, r, saver, log, req.Alias)
		}
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
		log.Error("failed to add url", sl.ErrorAttr(err))
		render.JSON(w, r, Api.Error("failed to add url"))
		return
	}

	log.Info("url added")

	render.JSON(w, r, response{
		Response: Api.Ok(),
		Alias:    alias,
	})
}

func saveWithAlias(w http.ResponseWriter, r *http.Request, saver urlSaver, log *slog.Logger, url string) {
	// todo remove literal
	aliasLen := 2
	alias := random.NewRandomString(aliasLen)

	for i := 0; i < 100; i++ {
		err := saver.SaveURL(url, alias)

		if i == 99 {
			log.Error("failed to generate url after 100 trys")
			render.JSON(w, r, Api.Error("failed to generate alias"))
			return
		}

		if errors.Is(err, storage.ErrAliasExists) {
			log.Info("existing alias was generated", slog.String("alias", alias), slog.Int("number", i))
			alias = random.NewRandomString(aliasLen)
			continue
		}

		if err != nil {
			log.Error("failed to add url", sl.ErrorAttr(err))
			render.JSON(w, r, Api.Error("failed to add url"))
			return
		}

		break
	}

	log.Info("url added", slog.String("alias", alias))

	render.JSON(w, r, response{
		Response: Api.Ok(),
		Alias:    alias,
	})
}
