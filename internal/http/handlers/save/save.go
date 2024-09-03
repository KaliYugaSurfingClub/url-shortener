package save

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"link_shortener/internal/http/middlewares/mwLogger"
	"link_shortener/internal/lib/Api"
	"link_shortener/internal/lib/sl"
	"link_shortener/internal/storage"
	"log/slog"
	"net/http"
	"net/url"
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

		if resp, err := validateRequest(req, r.Host); err != nil {
			log.Error("validation failed", sl.ErrorAttr(err))
			render.JSON(w, r, resp)
			return
		}

		alias, err := saver.SaveURL(req.URL, req.Alias, 1*time.Second)

		switch {
		//524
		case errors.Is(err, storage.NotEnoughTimeToGenerate):
			log.Info(err.Error())
			render.JSON(w, r, Api.Error(err.Error()))
			return
		//409
		case errors.Is(err, storage.ErrAliasExists):
			log.Info(err.Error())
			render.JSON(w, r, Api.Error(err.Error()))
			return
		//500
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

func validateRequest(req request, currentHost string) (Api.Response, error) {
	const op = "handlers.save.validateRequest"

	if errs := validator.New().Struct(req); errs != nil {
		return Api.ValidationError(errs.(validator.ValidationErrors)), errs
	}

	urlForShort, err := url.Parse(req.URL)
	if err != nil {
		return Api.Error("url is not correct"), fmt.Errorf("%s: %w", op, err)
	}

	//todo getBannedHosts is not implemented
	bannedHosts, err := getBannedHosts()
	if err != nil {
		return Api.Error("server error"), fmt.Errorf("%s: %w", op, err)
	}

	if _, contains := bannedHosts[urlForShort.Host]; contains || urlForShort.Host == currentHost {
		return Api.Error("deprecated url"), fmt.Errorf("%s: %w", op, err)
	}

	return Api.Ok(), nil
}

// todo move this from here
// maybe other service with noSQL db for adding new deprecated urls
func getBannedHosts() (map[string]struct{}, error) {
	const op = "getBannedHosts"

	res := make(map[string]struct{})

	////..adding deprecated ports to res....
	//
	//currentPort, err := http.
	//if err != nil {
	//	return nil, fmt.Errorf("%s: %w", op, err)
	//}
	//res[currentPort] = struct{}{}

	return res, nil
}
