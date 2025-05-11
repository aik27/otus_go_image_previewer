package routes

import (
	"context"
	"net/http"

	lrucache "github.com/aik27/otus_go_image_previewer/internal/cache/lru"
	"github.com/aik27/otus_go_image_previewer/internal/config"
	"github.com/aik27/otus_go_image_previewer/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ChiRouter(ctx context.Context, cfg *config.Config, cache lrucache.Cache) *chi.Mux {
	rootRouter := chi.NewRouter()

	rootRouter.Use(middleware.StripSlashes)

	rootRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
	})

	rootRouter.Get("/fill/{width}/{height}/*", func(w http.ResponseWriter, r *http.Request) {
		handlers.FillHandler(ctx, cfg, cache, w, r)
	})

	return rootRouter
}
