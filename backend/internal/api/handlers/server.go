package handlers

import (
	"github.com/go-chi/chi/v5"
)

func (h *OrderHandler) WithAuthHandlers(r chi.Router) {
	r.Route("/order", func(r chi.Router) {
		r.Get("/{orderUID}", h.GetOrder)
	})
}
