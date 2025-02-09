package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"money-manager/internal/database"
	"money-manager/internal/lib/mail/sender"
)

const label = "user"

type Handler struct {
	logger   *slog.Logger
	query    *database.Queries
	validate *validator.Validate
	mailer   sender.Sender
}

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
	})
}

func NewUserHandler(log *slog.Logger, db *database.Queries, mailer sender.Sender) *Handler {
	return &Handler{
		logger:   log,
		query:    db,
		validate: validator.New(),
		mailer:   mailer,
	}
}
