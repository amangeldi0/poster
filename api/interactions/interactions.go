package interactions

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"log/slog"
	authmiddleware "money-manager/api/middlewares/auth"
	"money-manager/internal/database"
)

type Handler struct {
	logger   *slog.Logger
	query    *database.Queries
	validate *validator.Validate
}

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/interactions", func(r chi.Router) {
		r.With(authmiddleware.JWTAuthRequired).Post("/like", handler.LikeEntity)
		r.With(authmiddleware.JWTAuthRequired).Post("/unlike", handler.UnlikeEntity)
		r.With(authmiddleware.JWTAuthRequired).Post("/comment", handler.Comment)
	})

}

func NewInteractionsHandlers(log *slog.Logger, db *database.Queries) *Handler {
	return &Handler{
		logger:   log,
		query:    db,
		validate: validator.New(),
	}
}
