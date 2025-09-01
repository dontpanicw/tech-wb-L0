package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"tech-wb-L0/backend/internal/service"
)

type OrderHandler struct {
	service service.OrderService
}

func NewHandler(service service.OrderService) (*OrderHandler, error) {
	return &OrderHandler{service: service}, nil
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderUID := chi.URLParam(r, "orderUID")
	if orderUID == "" {
		http.Error(w, "orderUID is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(ctx, orderUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get order: %v", err), http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(order)

}
