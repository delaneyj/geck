package ecs

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(setupCtx context.Context, world *World, baseRouter chi.Router) error {
	baseRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Geck!"))
	})

	baseRouter.Route("/sparsesets", func(sparseSetsRouter chi.Router) {
		sparseSetsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, SparseSets!"))
		})

		sparseSetsRouter.Route("/name", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.nameComponents
				SparseSetView(ss).Render(r.Context(), w)
			})
		})

	})

	return nil
}
