package main

import (
	"context"
	"example/internal/routes"
	"example/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"

	"example/internal/api"
	_ "example/internal/dto"
	"example/internal/repositories"
	_ "example/internal/services"
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

	// Инициализация репозиториев
	repo := repositories.NewOrderPostgresRepository(pool)

	//userRepo := repositories.NewUserPostgresRepository(pool)

	if err != nil {
		panic(err)
	}

	ordersService := services.NewOrdersService(repo)

	createOrdersHandler := api.NewCreateOrderHandler(ordersService)

	// Настройка маршрутов через отдельную функцию
	router := routes.SetupRoutes(createOrdersHandler)

	log.Println("API server started on :8080")
	if err := http.ListenAndServe(listenAddress, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Default().Printf("api server started on %s\n", listenAddress)
	if err := http.ListenAndServe(listenAddress, router); err != nil {
		panic(err)
	}
}
