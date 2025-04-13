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
	//var req struct {
	//	HotelID    string    `json:"hotel_id"`
	//	RoomTypeID string    `json:"room_type_id"`
	//	From       time.Time `json:"from"`
	//	To         time.Time `json:"to"`
	//	UserID     string    `json:"user_id"`
	//}

	orderRequest, err := parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validateCreateOrder(orderRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	order, err := h.ordersService.Create(&orderRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = orderCreatedResponse(w, order)
	if err != nil {
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
