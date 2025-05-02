package handlers

import "github.com/go-playground/validator/v10"

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text" validate:"required,min=3,max=500"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
