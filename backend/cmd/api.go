package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/zipshell/dev-learning-tracker/internal/adapters/postgresql/sqlc"
	"github.com/zipshell/dev-learning-tracker/internal/auth"
	"github.com/zipshell/dev-learning-tracker/internal/entries"
	"github.com/zipshell/dev-learning-tracker/internal/folders"
	mw "github.com/zipshell/dev-learning-tracker/internal/middleware"
	"github.com/zipshell/dev-learning-tracker/internal/users"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	})

	userService := users.NewUserService(app.db)
	userHandler := users.NewUserHandler(userService)
	r.Post("/users", userHandler.CreateUser)

	authService := auth.NewAuthService(repo.New(app.db))
	authHandler := auth.NewAuthHandler(authService)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)
	})

	folderService := folders.NewService(app.db)
	folderHandler := folders.NewHandler(folderService)
	r.Route("/folders", func(r chi.Router) {
		r.Use(mw.Auth(authService))
		r.Post("/", folderHandler.CreateFolder)
		r.Get("/{user_id}", folderHandler.ListFoldersByUserId)
	})

	entryService := entries.NewEntryService(repo.New(app.db))
	entryHandler := entries.NewEntryHandler(entryService)
	r.Route("/entries", func(r chi.Router) {
		r.Use(mw.Auth(authService))
		r.Post("/", entryHandler.CreateEntry)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at address %s", app.config.addr)

	return srv.ListenAndServe()
}
