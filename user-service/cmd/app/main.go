package main

import (
	"flower-shop/user-service/internal/config"
	del "flower-shop/user-service/internal/http-server/handlers/delete"
	"flower-shop/user-service/internal/http-server/handlers/login"
	"flower-shop/user-service/internal/http-server/handlers/signup"
	mw "flower-shop/user-service/internal/http-server/middleware"
	"flower-shop/user-service/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	log.Info("starting user-service")

	storage, err := postgres.New(cfg.StorageConn)
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("user-service/web/static/jss"))
	router.Handle("/jss/*", http.StripPrefix("/jss", fs))

	router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return mw.AuthMiddleware(next, cfg.Jwt)
		})

		r.Get("/delete", del.DeletePage)
		r.Post("/delete", del.New(log, storage, cfg.Jwt))
	})

	router.Get("/login", login.LoginPage)
	router.Post("/login", login.New(log, storage, cfg.Jwt))
	router.Get("/signup", signup.SignUpPage)
	router.Post("/signup", signup.New(log, storage))

	router.Group(func(r chi.Router) {})

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	log.Info("server started on port: ", slog.String("address", cfg.HttpServer.Address))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server is stopped")
}

func newLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}
