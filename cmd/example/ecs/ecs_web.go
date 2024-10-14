package ecs

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(setupCtx context.Context, world *World, baseRouter chi.Router) error {
	baseRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	baseRouter.Route("/sparsesets", func(sparseSetsRouter chi.Router) {
		sparseSetsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			AllSparseSetsView().Render(r.Context(), w)
		})

		sparseSetsRouter.Route("/names", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.nameComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/positions", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.positionComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/velocities", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.velocityComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/rotations", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.rotationComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/directions", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.directionComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/enemy", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.enemyTags
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/gravities", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.gravityComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/spaceship", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.spaceshipTags
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/spacestation", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.spacestationTags
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/factions", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.factionComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/docked_tos", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.dockedToComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/planet", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.planetTags
				SparseSetView(ss).Render(r.Context(), w)
			})

		})

		sparseSetsRouter.Route("/ruled_bys", func(ssRouter chi.Router) {
			ssRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ss := world.ruledByComponents
				SparseSetView(ss).Render(r.Context(), w)
			})

		})
	})

	return nil
}
