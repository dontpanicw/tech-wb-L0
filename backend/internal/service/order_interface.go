package service

import (
	"context"
	"tech-wb-L0/backend/domain"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrder(ctx context.Context, id string) (*domain.Order, error)
}
