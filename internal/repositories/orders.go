package repositories

import (
	"context"
	"errors"
	"example/internal/dto"
	"example/internal/interfaces"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrdersPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewOrderPostgresRepository(pool *pgxpool.Pool) interfaces.OrdersRepository {
	return &OrdersPostgresRepository{pool: pool}
}

// Create - создание нового заказа
func (r *OrdersPostgresRepository) Create(ctx context.Context, order *dto.Order) error {
	query := `
		INSERT INTO orders (
			id, user_id, hotel_id, room_type_id, 
			from_date, to_date, status
		) VALUES (
			$1, $2, $3, $4, 
			$5, $6, $7
		)`

	_, err := r.pool.Exec(ctx, query,
		order.UserID,
		order.HotelID,
		order.RoomTypeID,
		order.From,
		order.To,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (r *OrdersPostgresRepository) GetByID(ctx context.Context, id string) (*dto.Order, error) {
	query := `
		SELECT 
			id, user_id, hotel_id, room_type_id,
			from_date, to_date, status,
			created_at, updated_at
		FROM orders 
		WHERE id = $1`

	var order dto.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.UserID,
		&order.HotelID,
		&order.RoomTypeID,
		&order.From,
		&order.To,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// Update - обновление заказа
func (r *OrdersPostgresRepository) Update(ctx context.Context, order *dto.Order) error {
	query := `
		UPDATE orders SET
			user_id = $2,
			hotel_id = $3,
			room_type_id = $4,
			from_date = $5,
			to_date = $6,
			status = $7,
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query,
		order.UserID,
		order.HotelID,
		order.RoomTypeID,
		order.From,
		order.To,
	)

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// Delete - удаление заказа
func (r *OrdersPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}
