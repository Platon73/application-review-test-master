package api

import (
	"encoding/json"
	"errors"
	"example/internal/dto"
	"example/internal/interfaces"
	"log"
	"net/http"
)

type CreateOrderHandler struct {
	ordersService interfaces.OrdersService
}

func NewCreateOrderHandler(
	ordersService interfaces.OrdersService,
) *CreateOrderHandler {
	return &CreateOrderHandler{ordersService: ordersService}
}

func (h *CreateOrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	orderRequest, err := parseRequest(r)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCreateOrder(orderRequest); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.ordersService.Create(r.Context(), &orderRequest)
	if err != nil {
		http.Error(w, "failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := orderCreatedResponse(w, order); err != nil {
		log.Default().Print(err)
	}
}

func validateCreateOrder(order dto.Order) error {
	if order.From.After(order.To) {
		return errors.New("from after to")
	}
	return nil
}

func parseRequest(r *http.Request) (dto.Order, error) {
	var newOrder dto.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		return dto.Order{}, err
	}
	return newOrder, nil
}

func orderCreatedResponse(w http.ResponseWriter, order *dto.Order) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(order)
	if err != nil {
		return err
	}
	return nil
}
