package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adzil/bebek"
	"github.com/bukalapak/apinizer/response"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Handler is the room booking HTTP handler.
type Handler struct {
	Logger  *zap.Logger
	Service bebek.Service
}

// ArrayMeta is response meta data for array of object
type ArrayMeta struct {
	HTTPStatus int `json:"http_status"`
	Length     int `json:"length"`
}

// GetRooms serves GET /rooms endpoint.
func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	disableCache(w)

	rooms, err := h.Service.GetRooms()
	if err != nil {
		h.Logger.Error("get rooms", zap.Error(err))
		writeResponseError(w, err)
	} else {
		meta := newArrayMeta(len(rooms))
		sucResp := response.BuildSuccess(rooms, "", meta)
		response.Write(w, sucResp, http.StatusOK)
	}
}

// GetBookings serves GET /bookings endpoint.
func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	disableCache(w)

	req, errResp := getReversationRequestFromQuery(r)
	if errResp != nil {
		h.Logger.Error("get reservation request from query", zap.Error(errResp))
		response.Write(w, errResp, http.StatusBadRequest)
		return
	}
	resv, err := h.Service.GetReservations(*req)
	if err != nil {
		h.Logger.Error("get reservations", zap.Error(err))
		writeResponseError(w, err)
	} else {
		meta := newArrayMeta(len(resv))
		sucResp := response.BuildSuccess(resv, "", meta)

		response.Write(w, sucResp, http.StatusOK)
	}
}

// PostBooking serves POST /bookings endpoint.
func (h *Handler) PostBooking(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	disableCache(w)

	bookingReq, errResp := createBookingRequestFromBody(r)
	if errResp != nil {
		h.Logger.Error("create booking request from body", zap.Error(errResp))
		response.Write(w, errResp, http.StatusBadRequest)
		return
	}
	if valid := validateCreateBookingRequest(bookingReq); !valid {
		errResp := buildErrorResponse("body", "POST body is invalid")
		h.Logger.Error("validate create booking request", zap.Error(errResp))
		response.Write(w, errResp, http.StatusBadRequest)
		return
	}
	bookingID, err := h.Service.CreateBooking(*bookingReq)
	if err != nil {
		h.Logger.Error("create booking", zap.Error(err))
		writeResponseError(w, err)
		return
	}
	meta := struct {
		HTTPStatus int `json:"http_status"`
	}{
		HTTPStatus: http.StatusOK,
	}
	resp := struct {
		BookingID string `json:"booking_id"`
	}{
		BookingID: bookingID,
	}
	sucResp := response.BuildSuccess(resp, "", meta)
	response.Write(w, sucResp, http.StatusOK)
}

// DeleteBooking serves DELETE /bookings endpoint.
func (h *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	disableCache(w)

	bookingID := params.ByName("id")
	actor := r.Header.Get("X-Telegram-User")
	if actor == "" {
		errResp := buildErrorResponse("header", "User is not set")
		h.Logger.Error("get user header", zap.Error(errResp))
		response.Write(w, errResp, http.StatusBadRequest)
		return
	}
	req := bebek.DeleteBookingRequest{
		Actor:     actor,
		BookingID: bookingID,
	}
	err := h.Service.DeleteBooking(req)
	if err != nil {
		h.Logger.Error("delete booking", zap.Error(err))
		writeResponseError(w, err)
		return
	}
	meta := struct {
		HTTPStatus int `json:"http_status"`
	}{
		HTTPStatus: http.StatusOK,
	}
	sucResp := response.BuildSuccess(nil, "Delete Booking successful", meta)
	response.Write(w, sucResp, http.StatusOK)
}

// GetSelfBookings serves GET /self/bookings endpoint.
func (h *Handler) GetSelfBookings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	disableCache(w)

	actor := r.Header.Get("X-Telegram-User")
	if actor == "" {
		errResp := buildErrorResponse("header", "User is not set")
		h.Logger.Error("get user header", zap.Error(errResp))
		response.Write(w, errResp, http.StatusBadRequest)
		return
	}
	req := bebek.GetSelfReservationsRequest{
		Actor: actor,
	}
	resv, err := h.Service.GetSelfReservations(req)
	if err != nil {
		h.Logger.Error("create booking", zap.Error(err))
		writeResponseError(w, err)
		return
	}
	meta := newArrayMeta(len(resv))
	sucResp := response.BuildSuccess(resv, "", meta)
	response.Write(w, sucResp, http.StatusOK)
}

func buildErrorResponse(field, msg string) response.ErrorBody {
	return response.BuildError(response.ErrorInfo{
		Code:    http.StatusBadRequest,
		Field:   field,
		Message: msg,
	})
}

func createBookingRequestFromBody(r *http.Request) (*bebek.CreateBookingRequest, *response.ErrorBody) {
	if r.Body == nil {
		errResp := buildErrorResponse("body", "POST body is empty")
		return nil, &errResp
	}
	var bookingReq bebek.CreateBookingRequest
	err := json.NewDecoder(r.Body).Decode(&bookingReq)
	if err != nil {
		errResp := buildErrorResponse("body", err.Error())
		return nil, &errResp
	}
	actor := r.Header.Get("X-Telegram-User")
	if actor == "" {
		errResp := buildErrorResponse("header", "User is not set")
		return nil, &errResp
	}
	bookingReq.Actor = actor
	return &bookingReq, nil
}

func getReversationRequestFromQuery(r *http.Request) (*bebek.GetReservationsRequest, *response.ErrorBody) {
	query := r.URL.Query()
	roomID := query.Get("room_id")
	dateStr := query.Get("date")
	date := time.Now()
	if dateStr != "" {
		var err error
		date, err = time.Parse(bebek.DateLayout, dateStr)
		if err != nil {
			errResp := buildErrorResponse("date", "date format is invalid")
			return nil, &errResp
		}
	}
	return &bebek.GetReservationsRequest{
		Date:   bebek.Date{Time: date},
		RoomID: roomID,
	}, nil
}

func newArrayMeta(length int) *ArrayMeta {
	return &ArrayMeta{
		HTTPStatus: http.StatusOK,
		Length:     length,
	}
}

func validateCreateBookingRequest(req *bebek.CreateBookingRequest) bool {
	if req.RoomID == "" {
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	if req.Date.Before(today) {
		return false
	}
	if req.End < req.Start {
		return false
	}
	return true
}

func writeResponseError(w http.ResponseWriter, err error) {
	errResp := response.BuildError(err)
	response.Write(w, errResp, http.StatusInternalServerError)
}

func disableCache(w http.ResponseWriter) {
	header := w.Header()
	header.Add("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Add("Pragma", "no-cache")
	header.Add("Expires", "0")
}
