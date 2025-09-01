package repository

import (
	"context"
	"tech-wb-L0/backend/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrder(ctx context.Context, orderUID string) (*domain.Order, error)
}
