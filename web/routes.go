package web

import (
	"net/http"

	"github.com/delaneyj/geck/cmd/example/ecs"
	"github.com/go-chi/chi/v5"
)

func SparseSetRoutes[T any](baseRouter chi.Router, ss *ecs.SparseSet[T]) {

	baseRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
	})
}
