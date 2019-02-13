package bebek

import "time"

type GetReservationsRequest struct {
	RoomID string
	Date   time.Time
}

type GetSelfReservationsRequest struct {
	Actor string
}

type CreateBookingRequest struct {
	Booking *Booking
}

type DeleteBookingRequest struct {
	BookingID string
	Actor     string
}

type Service interface {
	GetRooms() ([]*Room, error)
	GetReservations(GetReservationsRequest) ([]*Reservation, error)
	GetSelfReservations(GetSelfReservationsRequest) ([]*Reservation, error)
	CreateBooking(CreateBookingRequest) error
	DeleteBooking(DeleteBookingRequest) error
}
