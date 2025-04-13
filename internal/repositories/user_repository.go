package repositories

import (
	"context"
	"errors"
	"example/internal/dto"
	"example/internal/interfaces"
	"example/internal/utils"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserPostgresRepository struct {
	Pool *pgxpool.Pool
}

func NewUserPostgresRepository(pool *pgxpool.Pool) interfaces.UserRepository {
	return &UserPostgresRepository{Pool: pool}
}

func (r *UserPostgresRepository) GetByID(ctx context.Context, id string) (*dto.User, error) {
	query := `
		SELECT id, firstname, lastname, email, age, created
		FROM users
		WHERE id = $1
	`

	var user dto.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Age,
		&user.Created,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserPostgresRepository) Create(user *dto.User) error {
	query := `
		INSERT INTO users (id, firstname, lastname, email, age, created)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.Pool.Exec(context.Background(), query,
		user.ID,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Age,
		user.Created,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserPostgresRepository) Update(user *dto.User) error {
	query := `
		UPDATE users
		SET 
			firstname = $2,
			lastname = $3,
			email = $4,
			age = $5,
			created = $6
		WHERE id = $1
	`

	result, err := r.Pool.Exec(context.Background(), query,
		user.ID,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Age,
		user.Created,
	)

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
