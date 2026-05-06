package update

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Alias  string `json:"alias" validate:"required"`
	NewURL string `json:"new_url" validate:"required,url"`
}

type Response struct {
	resp.Response
}

type UpdateURL interface {
	UpdateURL(alias string, newURL string) (int64, error)
}

func New(log *slog.Logger, UpdateURL UpdateURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const update = "handlers.url.update.New"

		log = log.With(
			slog.String("fn", update),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("invalid request body", "error", err)
			render.JSON(w, r, resp.Error("invalid request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request body", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		newURL := req.NewURL

		id, err := UpdateURL.UpdateURL(alias, newURL)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to update url", "alias", alias, sl.Err(err))
			render.JSON(w, r, resp.Error("failed to update url"))
			return
		}

		log.Info("url updated", slog.String("alias", alias), slog.Int64("id", id))

		responseOK(w, r)

	}

}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
