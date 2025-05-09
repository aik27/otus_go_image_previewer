package routes

import (
	"net/http"

	"github.com/aik27/otus_go_image_previewer/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func ChiRouter() *chi.Mux {
	rootRouter := chi.NewRouter()

	rootRouter.Use(middleware.StripSlashes)

	rootRouter.Get("/health", healthHandler)

	return rootRouter
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response.OK())
}
