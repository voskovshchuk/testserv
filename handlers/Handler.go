package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	messages = make(map[int]Message)
	nextID   = 1
)

func GetHandler(c echo.Context) error {
	c.Logger().Info("Handling GET request for all messages")

	messageList := make([]Message, 0, len(messages))
	for _, msg := range messages {
		messageList = append(messageList, msg)
	}

	c.Logger().Infof("Returning %d messages", len(messageList))
	return c.JSON(http.StatusOK, messageList)
}

func PostHandler(c echo.Context) error {
	c.Logger().Info("Handling POST request for new message")

	var msg Message
	if err := c.Bind(&msg); err != nil {
		c.Logger().Error("Failed to bind JSON:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&msg); err != nil {
		c.Logger().Error("Validation failed:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	msg.ID = nextID
	nextID++
	messages[msg.ID] = msg

	c.Logger().Infof("Added new message with ID: %d", msg.ID)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message added successfully",
	})
}

func DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	c.Logger().Infof("Handling DELETE request for message ID: %s", idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Logger().Error("Invalid ID format:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid message ID format",
		})
	}

	if _, exists := messages[id]; !exists {
		c.Logger().Warnf("Message with ID %d not found", id)
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Message not found",
		})
	}

	delete(messages, id)
	c.Logger().Infof("Deleted message with ID: %d", id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message deleted successfully",
	})
}

func PatchHandler(c echo.Context) error {
	idParam := c.Param("id")
	c.Logger().Infof("Handling PATCH request for message ID: %s", idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Logger().Error("Invalid ID format:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid message ID format",
		})
	}

	if _, exists := messages[id]; !exists {
		c.Logger().Warnf("Message with ID %d not found", id)
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Message not found",
		})
	}

	var update Message
	if err := c.Bind(&update); err != nil {
		c.Logger().Error("Failed to bind JSON:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&update); err != nil {
		c.Logger().Error("Validation failed:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	update.ID = id
	messages[id] = update

	c.Logger().Infof("Updated message with ID: %d", id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message updated successfully",
	})
}
