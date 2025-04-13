package main

import (
	"context"
	"example/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"

	"example/internal/api"
	_ "example/internal/dto"
	"example/internal/repositories"
	_ "example/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const listenAddress = ":8080"

func main() {
	connStr := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	// Парсинг конфигурации пула соединений
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Настройки пула
	config.MaxConns = 50
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.HealthCheckPeriod = 30 * time.Second

	// Создание пула
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}
	defer pool.Close()

	repo := repositories.NewPostgresRepository(pool)
	// Инициализация репозиториев
	orderRepo := repositories.NewOrderPostgresRepository(pool)
	userRepo := repositories.NewUserPostgresRepository(pool)

	// Инициализация сервиса
	ordersService := services.NewOrdersService(
		orderRepo,
		userRepo, // Добавлен новый репозиторий
		emailSender,
	)
	if err != nil {
		panic(err)
	}

	// Инициализация предопределенных номеров
	_, err = repo.(*repositories.PostgresRepository).Pool.Exec(context.Background(), `
        INSERT INTO rooms(hotel_id, room_type_id) 
        VALUES 
            ('reddison', 'lux'),
            ('reddison', 'premium') 
        ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Initialization error: %v", err)
	}

	ordersService := services.NewOrdersService(repo)

	createOrdersHandler := api.NewCreateOrderHandler(ordersService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/orders", createOrdersHandler.Handle)

	log.Default().Printf("api server started on %s\n", listenAddress)
	if err := http.ListenAndServe(listenAddress, router); err != nil {
		panic(err)
	}
}
