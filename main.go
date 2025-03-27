package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Text string `json:"text"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var messages []Message

func GetHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &messages)
}

func PostHandler(c echo.Context) error {
	var message Message
	// Bind переводит Json в сообщение
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Не смогли добавить сообщение",
		})
	}
	messages = append(messages, message)
	return c.JSON(http.StatusOK, Response{
		Status:  "Dobavil",
		Message: "Message dobavlen",
	})
}
func DeleteHandler(c echo.Context) error {
	if len(messages) == 0 {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "No messages to delete",
		})
	}
	lastmessage := messages[len(messages)-1]
	messages = messages[:len(messages)-1]

	return c.JSON(http.StatusOK, Response{
		Status:  "Vse ok",
		Message: "Ybral message:" + lastmessage.Text,
	})
}

func main() {
	e := echo.New()
	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHandler)
	e.DELETE("/messages", DeleteHandler)

	e.Start("localhost:8080")
}
