package routes

import (
	"net/http"

	"example/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(createOrdersHandler *api.CreateOrderHandler) http.Handler {
	router := chi.NewRouter()

	// Подключение middleware
	router.Use(middleware.Logger)

	// Настройка маршрутов
	router.Post("/orders", createOrdersHandler.Handle)

	return router
}
