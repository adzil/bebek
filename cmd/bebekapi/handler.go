package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler is the room booking HTTP handler.
type Handler struct{}

// GetRooms serves GET /rooms endpoint.
func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

// GetBookings serves GET /bookings endpoint.
func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

// PostBooking serves POST /bookings endpoint.
func (h *Handler) PostBooking(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

// DeleteBooking serves DELETE /bookings endpoint.
func (h *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

// GetSelfBookings serves GET /self/bookings endpoint.
func (h *Handler) GetSelfBookings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
