package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

var (
	messages = make(map[int]Message)
	nextID   = 1
	db       *pgx.Conn
)

func SetDB(conn *pgx.Conn) {
	db = conn
}

func GetHandler(c echo.Context) error {
	c.Logger().Info("Handling GET request for all messages")

	// Запрос всех сообщений из БД
	rows, err := db.Query(context.Background(), "SELECT id, text, created_at FROM messages ORDER BY created_at DESC")

	if err != nil {
		c.Logger().Error("Database query error:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Failed to fetch messages from database",
		})
	}

	defer rows.Close()

	//запись из бд для передачи
	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Text, &msg.CreatedAt); err != nil {
			c.Logger().Error("Failed to scan row:", err)
			continue // Пропускаем битые строки, но можно и вернуть ошибку
		}
		messages = append(messages, msg)
	}

	c.Logger().Infof("Returning %d messages from database", len(messages))
	return c.JSON(http.StatusOK, messages)
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

	_, err := db.Exec(context.Background(),
		"INSERT INTO messages (text, created_at) VALUES ($1, $2)",
		msg.Text, time.Now())

	if err != nil {
		c.Logger().Error("Database insert error:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Failed to save message to database",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message added successfully",
	})
}

func DeleteHandler(c echo.Context) error {
	// Получаем ID сообщения из URL-параметра
	idParam := c.Param("id")
	c.Logger().Infof("Handling DELETE request for message ID: %s", idParam)

	// Преобразуем ID в число
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Logger().Error("Invalid ID format:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid message ID format",
		})
	}

	// Выполняем запрос на удаление в БД
	_, err = db.Exec(context.Background(),
		"DELETE FROM messages WHERE id = $1",
		id)

	if err != nil {
		c.Logger().Error("Database delete error:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Failed to delete message from database",
		})
	}

	// Проверяем, была ли удалена хотя бы одна строка

	c.Logger().Infof("Successfully deleted message with ID: %d", id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message deleted successfully",
	})
}

func PatchHandler(c echo.Context) error {
	// Получаем ID сообщения из URL-параметра
	idParam := c.Param("id")
	c.Logger().Infof("Handling PATCH request for message ID: %s", idParam)

	// Преобразуем ID в число
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Logger().Error("Invalid ID format:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid message ID format",
		})
	}

	// Парсим тело запроса
	var update Message
	if err := c.Bind(&update); err != nil {
		c.Logger().Error("Failed to bind JSON:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid request body",
		})
	}

	// Валидация данных
	if err := c.Validate(&update); err != nil {
		c.Logger().Error("Validation failed:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: err.Error(),
		})
	}

	// Обновляем сообщение в базе данных
	_, err = db.Exec(
		c.Request().Context(),
		"UPDATE messages SET text = $1 WHERE id = $2",
		update.Text,
		id,
	)

	if err != nil {
		c.Logger().Error("Database update error:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "Error",
			Message: "Failed to update message in database",
		})
	}

	c.Logger().Infof("Successfully updated message with ID: %d", id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message updated successfully",
		Data:    map[string]int{"id": id},
	})
}
