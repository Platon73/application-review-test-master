package dto

import "time"

type User struct {
	ID        string    `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Age       uint      `json:"age"`
	Created   time.Time `json:"created"`
}
