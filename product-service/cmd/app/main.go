package main

import (
	"flower-shop/product-service/internal/config"
	"flower-shop/product-service/internal/http-server/handlers/create"
	del "flower-shop/product-service/internal/http-server/handlers/delete"
	"flower-shop/product-service/internal/http-server/handlers/get"
	"flower-shop/product-service/internal/http-server/handlers/list"
	mw "flower-shop/product-service/internal/http-server/middleware"
	"flower-shop/product-service/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	log.Info("starting product-service")

	storage, err := postgres.New(cfg.StorageConn)
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("product-service/web/static/js"))
	router.Handle("/js/*", http.StripPrefix("/js", fs))

	router.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return mw.AuthMiddleware(next, cfg.Jwt)
		})

		r.Get("/create", create.CreatePage)
		r.Post("/create", create.New(log, storage, cfg.Jwt))
		r.Get("/product/{id}/page", get.ProductPage)
		r.Get("/product/{id}", get.New(log, storage))
		r.Delete("/product/{id}", del.New(log, storage, cfg.Jwt))
	})

	router.Get("/", list.ProductsPage)
	router.Get("/products", list.New(log, storage))

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
