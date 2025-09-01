package order

import (
	"context"
	"tech-wb-L0/backend/domain"
	"tech-wb-L0/backend/internal/repository"
)

type Order struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) (*Order, error) {
	return &Order{
		repo: repo,
	}, nil
}

func (o *Order) CreateOrder(ctx context.Context, order *domain.Order) error {
	return o.repo.CreateOrder(ctx, order)
}

func (o *Order) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return o.repo.GetOrder(ctx, id)
}
