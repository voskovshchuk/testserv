package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var messages = make(map[int]Message)
var nextID = 1

func GetHandler(c echo.Context) error {
	var messageMap []Message
	for _, msg := range messages {
		messageMap = append(messageMap, msg)
	}
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
	message.ID = nextID
	nextID++

	messages[message.ID] = message
	return c.JSON(http.StatusOK, Response{
		Status:  "Dobavil",
		Message: "Message dobavlen",
	})
}

func DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Не верный айди",
		})
	}
	if _, i := messages[id]; !i {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Нет сообщения",
		})
	}
	delete(messages, id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Сообщение удалено",
	})
}

func PatchHandler(c echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Не верный айди",
		})
	}

	var updateMessage Message
	// Bind переводит Json в сообщение
	if err := c.Bind(&updateMessage); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Не смогли обновить сообщение",
		})
	}

	if _, i := messages[id]; !i {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Нет сообщения",
		})
	}

	updateMessage.ID = id
	messages[id] = updateMessage

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Сообщение обновлено",
	})
}

func main() {

	e := echo.New()
	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHandler)
	e.DELETE("/messages/:id", DeleteHandler)
	e.PATCH("/messages/:id", PatchHandler)

	e.Start("localhost:8080")

}
