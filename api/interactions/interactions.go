package interactions

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	authmiddleware "poster/api/middlewares/auth"
	"poster/internal/database"
)

type Handler struct {
	logger   *slog.Logger
	query    *database.Queries
	validate *validator.Validate
}

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/interactions", func(r chi.Router) {
		r.Route("/post", func(r chi.Router) {
			r.With(authmiddleware.JWTAuthRequired).Post("/like/{id}", handler.LikePost)
			r.With(authmiddleware.JWTAuthRequired).Post("/unlike/{id}", handler.UnlikePost)
		})

		r.Route("/comment", func(r chi.Router) {
			r.With(authmiddleware.JWTAuthRequired).Post("/", handler.Comment)
			r.With(authmiddleware.JWTAuthRequired).Put("/{id}", handler.UpdateComment)
			r.With(authmiddleware.JWTAuthRequired).Delete("/{id}", handler.DeleteComment)

			r.With(authmiddleware.JWTAuthRequired).Post("/like/{id}", handler.LikeComment)
			r.With(authmiddleware.JWTAuthRequired).Post("/unlike/{id}", handler.UnlikeComment)
		})
	})

}

func NewInteractionsHandlers(log *slog.Logger, db *database.Queries) *Handler {
	return &Handler{
		logger:   log,
		query:    db,
		validate: validator.New(),
	}
}

func (h *Handler) isCommentExist(ctx context.Context, id uuid.UUID) bool {
	_, err := h.query.GetComment(ctx, id)
	return err == nil
}

func (h *Handler) isPostExist(ctx context.Context, id uuid.UUID) bool {
	_, err := h.query.GetPost(ctx, id)
	return err == nil
}
