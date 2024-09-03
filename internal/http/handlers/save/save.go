package save

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"link_shortener/internal/http/middlewares/mwLogger"
	"link_shortener/internal/lib/Api"
	"link_shortener/internal/lib/sl"
	"link_shortener/internal/storage"
	"log/slog"
	"net/http"
	"time"
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
	SaveURL(urlToSave string, alias string, timeToGenerate time.Duration) (string, error)
}

func New(saver urlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := mwLogger.GetCtxLog(r.Context(), "handlers.save.New")

		var req request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("cannot decode request", sl.ErrorAttr(err))
			render.JSON(w, r, Api.Error("cannot decode request"))
			return
		}

		log.Debug("request was decoded", slog.Any("request", req))

		if errs := validator.New().Struct(req); errs != nil {
			log.Error("validation error", slog.String("error", errs.Error()))

			errForUser := Api.ValidationError(errs.(validator.ValidationErrors))
			render.JSON(w, r, errForUser)
			return
		}

		alias, err := saver.SaveURL(req.URL, req.Alias, 1*time.Second)

		switch {
		case errors.Is(err, storage.NotEnoughTimeToGenerate):
			log.Info(err.Error())
			render.JSON(w, r, Api.Error(err.Error()))
			return
		//409
		case errors.Is(err, storage.ErrAliasExists):
			log.Info(err.Error())
			render.JSON(w, r, Api.Error(err.Error()))
			return
		case err != nil:
			log.Error("failed to add url", sl.ErrorAttr(err))
			render.JSON(w, r, Api.Error("failed to add url"))
			return
		}

		log.Info("alias was added", slog.String("alias", alias))

		render.JSON(w, r, response{
			Response: Api.Ok(),
			Alias:    alias,
		})
	}
}
