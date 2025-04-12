package repositories

import (
	"context"
	"errors"
	_ "time"

	"example/internal/dto"
	"example/internal/interfaces"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(connString string) (interfaces.OrdersRepository, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &postgresRepository{pool: pool}, nil
}

func (p *postgresRepository) Create(order *dto.Order) error {
	ctx := context.Background()

	// Проверка существования комнаты
	var exists bool
	err := p.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM rooms WHERE hotel_id=$1 AND room_type_id=$2)",
		order.HotelID, order.RoomTypeID).Scan(&exists)

	if err != nil || !exists {
		return errors.New("unknown room")
	}

	// Проверка доступности
	var overlappingOrders int
	err = p.pool.QueryRow(ctx, `
        SELECT COUNT(*) FROM orders 
        WHERE hotel_id=$1 AND room_type_id=$2 
        AND ($3 < "to") AND ($4 > "from")`,
		order.HotelID, order.RoomTypeID, order.From, order.To).Scan(&overlappingOrders)

	if err != nil || overlappingOrders > 0 {
		return errors.New("room not available")
	}

	// Вставка брони
	_, err = p.pool.Exec(ctx, `
        INSERT INTO orders(hotel_id, room_type_id, "from", "to") 
        VALUES($1, $2, $3, $4)`,
		order.HotelID, order.RoomTypeID, order.From, order.To)

	return err
}
