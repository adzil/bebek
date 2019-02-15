package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/adzil/bebek"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		logger.Panic("open database", zap.Error(err))
	}

	repository := &bebek.MySQLRepository{
		DB: db,
	}
	service := &bebek.Bebek{
		Repository: repository,
	}
	handler := &Handler{
		Service: service,
		Logger:  logger,
	}

	router := httprouter.New()
	router.GET("/rooms", handler.GetRooms)
	router.GET("/bookings", handler.GetBookings)
	router.POST("/bookings", handler.PostBooking)
	router.DELETE("/bookings/:id", handler.DeleteBooking)
	router.GET("/self/bookings", handler.GetSelfBookings)

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server closed", zap.Error(err))
	}
}
