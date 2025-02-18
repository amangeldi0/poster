package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	slogchi "github.com/samber/slog-chi"
	"gopkg.in/gomail.v2"
	"log/slog"
	"money-manager/api/auth"
	"money-manager/api/interactions"
	"money-manager/api/posts"
	"money-manager/internal/config"
	"money-manager/internal/database"
	"money-manager/internal/lib/logger/prettylogger"
	"money-manager/internal/lib/logger/sl"
	"money-manager/internal/lib/mail/sender"
	"net/http"
	"os"
)

func main() {

	// Configuring

	cfg, err := config.New()

	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}

	logger := setupLogger(cfg.Env)
	logger.Info("Starting money manager")

	// Connecting to Database

	db, err := sql.Open("postgres", cfg.Database.Address)

	if err != nil {
		logger.Error("failed to connect to database", sl.Err(err))
	}

	queries := database.New(db)

	// Mailer
	dialer, err := cfg.Mailer.Dialer.Dial()

	if err != nil {
		slog.Error("Error connecting to mailer", slog.Any("error", err))
		os.Exit(1)
	}

	defer func(dialer gomail.SendCloser) {
		err = dialer.Close()
		if err != nil {
			slog.Error("Error connecting to mailer", slog.Any("error", err))
			os.Exit(1)
		}
	}(dialer)

	mailer := sender.NewSender(cfg.Mailer.Email, cfg.Mailer.Dialer)

	// Routes

	router := chi.NewRouter()
	router.Use(slogchi.New(logger))

	usersHandlers := auth.NewAuthHandler(logger, queries, mailer)
	auth.RegisterRoutes(router, usersHandlers)

	postsHandlers := posts.NewPostsHandler(logger, queries)
	posts.RegisterRoutes(router, postsHandlers)

	interactionsHandlers := interactions.NewInteractionsHandlers(logger, queries)
	interactions.RegisterRoutes(router, interactionsHandlers)

	// Serving

	logger.Info("✅ Server started", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}

}

func setupLogger(level string) *slog.Logger {

	var log *slog.Logger

	if level == "dev" {
		prettyHandler := prettylogger.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   false,
			ReplaceAttr: nil,
		})

		log = slog.New(prettyHandler)
	} else if level == "prod" {
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	} else {
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return log

}
