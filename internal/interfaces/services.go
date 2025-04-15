package interfaces

import (
	"context"
	"example/internal/dto"
)

type OrdersService interface {
	Create(ctx context.Context, order *dto.Order) (*dto.Order, error)
}
