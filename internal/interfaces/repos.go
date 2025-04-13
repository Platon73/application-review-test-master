package interfaces

import (
	"context"
	"example/internal/dto"
)

type OrdersRepository interface {
	Create(ctx context.Context, order *dto.Order) error
	GetByID(ctx context.Context, id string) (*dto.Order, error)
	Update(ctx context.Context, order *dto.Order) error
	Delete(ctx context.Context, id string) error
}
