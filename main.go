package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"serv/handlers"

	"github.com/jackc/pgx/v5"
)

func main() {
	e := echo.New()

	//Инициализация валидатора
	e.Validator = &handlers.CustomValidator{Validator: validator.New()}

	//логирование
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} ${method} ${uri} status: ${status}` + "\n",
	}))

	e.Use(middleware.Recover())

	e.GET("/messages", handlers.GetHandler)
	e.POST("/messages", handlers.PostHandler)
	e.DELETE("/messages/:id", handlers.DeleteHandler)
	e.PATCH("/messages/:id", handlers.PatchHandler)

	// соединение с бд
	connStr := "postgres://postgres:12345@127.0.0.1:5438/postgres"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	// Проверка соединения
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PostgreSQL!")

	handlers.SetDB(conn)

	// Запуск сервера
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Server failed: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Print("Shutdown error: ", err)
	}

}
