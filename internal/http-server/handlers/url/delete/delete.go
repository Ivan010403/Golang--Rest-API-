package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UrlDeleter interface {
	DeleteUrl(alias string) error
}

func New(log *slog.Logger, urlDeleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("not found for deleting"))
		}

		err := urlDeleter.DeleteUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, response.Error("not found for deleting"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("delete url", slog.String("url", alias))
	}
}
