package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	slogchi "github.com/samber/slog-chi"
	"log/slog"
	"money-manager/internal/config"
	"money-manager/internal/lib/logger/prettylogger"
	"money-manager/internal/lib/logger/sl"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.New()

	router := chi.NewRouter()

	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
	}

	logger := setupLogger(cfg.Env)
	logger.Info("Starting money manager")

	router.Use(slogchi.New(logger))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info("✅ Server started", slog.String("address", cfg.HTTPServer.Address))

	if err = srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server", sl.Err(err))
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
