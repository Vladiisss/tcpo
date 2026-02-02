package httpx

import (
	"github.com/go-chi/chi/v5"

	"example.com/notes-api/internal/http/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", h.ListNotes)
		r.Post("/", h.CreateNote)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetNote)
			r.Patch("/", h.PatchNote)
			r.Delete("/", h.DeleteNote)
		})
	})

	return r
}
