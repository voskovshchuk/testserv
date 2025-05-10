package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	ID        int       `json:"id" db:"id"`
	Text      string    `json:"text" validate:"required,min=3,max=500" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
