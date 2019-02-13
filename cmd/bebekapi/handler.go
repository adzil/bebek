package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct{}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (h *Handler) PostBooking(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (h *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}

func (h *Handler) GetSelfBookings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
