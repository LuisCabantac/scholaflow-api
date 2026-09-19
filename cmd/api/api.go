package main

import (
	"log/slog"
	"net/http"
	"time"

	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"
	"github.com/LuisCabantac/scholaflow-api/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
)

type dbConfig struct {
	dsn string
}

type config struct {
	addr    string
	env     string
	appUrl  string
	authUrl string
	db      dbConfig
}

type application struct {
	s3Client   *s3.Client
	mailClient *mail.Client
	db         *pgxpool.Pool
	queries    *postgres.Queries
	cfg        config
}

const apiName = "ScholaFlow API"

const version = "1.0.0"

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	allowedOrigins := []string{app.cfg.appUrl}
	if app.cfg.env == "development" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000", "http://localhost:5173")
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, struct {
			Message     string `json:"message"`
			Version     string `json:"version"`
			Status      string `json:"status"`
			Environment string `json:"environment"`
		}{
			Message:     apiName,
			Version:     version,
			Status:      "running",
			Environment: app.cfg.env,
		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.Text(w, http.StatusOK, "OK")
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.cfg.addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 35 * time.Second,
		IdleTimeout:  time.Minute,
	}

	slog.Info("starting server", "addr", app.cfg.addr, "env", app.cfg.env)

	return srv.ListenAndServe()
}
