package dto

import "time"

type Order struct {
	HotelID    string    `json:"hotel_id"`
	RoomTypeID string    `json:"room_id"`
	UserEmail  string    `json:"email"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	UserID     string    `json:"user_id"` // Связь через ID пользователя
	User       *User     `json:"user"`
}
