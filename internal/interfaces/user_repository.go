package interfaces

import (
	"context"
	"example/internal/dto"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*dto.User, error) // Добавить context
	Create(user *dto.User) error
	Update(user *dto.User) error
}
