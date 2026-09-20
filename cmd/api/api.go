package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"
	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
	"github.com/LuisCabantac/scholaflow-api/internal/classrooms"
	"github.com/LuisCabantac/scholaflow-api/internal/request"
	"github.com/LuisCabantac/scholaflow-api/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/lestrrat-go/jwx/v4/jwk"
	"github.com/lestrrat-go/jwx/v4/jwt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
)

const apiName = "ScholaFlow API"

const version = "1.0.0"

type dbConfig struct {
	dsn string
}

type jwksCache struct {
	mu          sync.RWMutex
	keySet      jwk.Set
	lastFetched time.Time
	ttl         time.Duration
}

type config struct {
	addr    string
	env     string
	appURL  string
	authURL string
	db      dbConfig
}

type application struct {
	s3Client   *s3.Client
	mailClient *mail.Client
	db         *pgxpool.Pool
	queries    *postgres.Queries
	jwks       *jwksCache
	cfg        config
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	allowedOrigins := []string{app.cfg.appURL}
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

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(app.authenticate)

		r.Route("/classrooms", func(r chi.Router) {
			classroomSvc := classrooms.NewService(app.queries)
			classroomHandler := classrooms.NewHandler(classroomSvc)

			r.Post("/", classroomHandler.Create)
			r.Post("/join", classroomHandler.Enroll)

			r.Route("/{classroomID}", func(r chi.Router) {
				r.Post("/leave", classroomHandler.Unenroll)
			})
		})

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

func (app *application) getCachedKeySet() jwk.Set {
	app.jwks.mu.RLock()
	defer app.jwks.mu.RUnlock()

	if app.jwks.keySet != nil && time.Since(app.jwks.lastFetched) < app.jwks.ttl {
		return app.jwks.keySet
	}

	return nil
}

func (app *application) getKeySet(ctx context.Context) (jwk.Set, error) {
	if set := app.getCachedKeySet(); set != nil {
		return set, nil
	}

	app.jwks.mu.Lock()
	defer app.jwks.mu.Unlock()

	if app.jwks.keySet != nil && time.Since(app.jwks.lastFetched) < app.jwks.ttl {
		return app.jwks.keySet, nil
	}

	jwksURL := fmt.Sprintf("%s/api/auth/jwks", app.cfg.authURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	set, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, err
	}

	app.jwks.keySet = set
	app.jwks.lastFetched = time.Now()

	return app.jwks.keySet, nil
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keySet, err := app.getKeySet(r.Context())
		if err != nil {
			response.Error(w, apperrors.ErrInternalServerError)
			return
		}

		token, err := jwt.ParseRequest(r, jwt.WithKeySet(keySet), jwt.WithValidate(true))
		if err != nil {
			response.Error(w, apperrors.ErrUnauthorizedAccess)
			return
		}
		userID, ok := token.Subject()
		if !ok || userID == "" {
			response.Error(w, apperrors.ErrUnauthorizedAccess)
			return
		}

		email, _ := jwt.Get[string](token, "email")
		name, _ := jwt.Get[string](token, "name")

		authUsr := &request.AuthUser{
			ID:    userID,
			Email: email,
			Name:  name,
		}

		next.ServeHTTP(w, request.SetAuthUser(r, authUsr))
	})
}
