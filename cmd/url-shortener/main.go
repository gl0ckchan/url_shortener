package main

import (
	"os"
	"net/http"
	"log/slog"

	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev" 
	envProd  = "prod"

	googleLink = "https://google.com/"
	googleAlias = "google"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	// log = log.With(slog.String("env", cfg.Env))
	log.Info("starting url-shortener")
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage:", slog.Any("error", err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/delete", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string {
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		
		r.Post("/{alias}",  delete.New(log, storage))
		r.Delete("/{alias}",  delete.New(log, storage))	
	})

	router.Post("/save",    save.New(log, storage))	
	router.Get ("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:    	  cfg.Address,
		Handler: 	  router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
		case envLocal:
			log = slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		case envDev:
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)
		case envProd:
			log = slog.New(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			)
	}

	return log
}
