package main

import (
	"net/http"

	"go.uber.org/zap"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger, _ := zap.NewProduction()

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		logger.Panic("open database", zap.Error(err))
	}

	_ = db

	handler := &Handler{}

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
