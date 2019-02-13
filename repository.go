package bebek

import "time"

type GetBookingsRequest struct {
	RoomID     string
	Date       time.Time
	StartFrom  bool
	ReservedBy string
}

type Repository interface {
	GetRooms() ([]*Room, error)
	GetBookings(GetBookingsRequest) ([]*Booking, error)
	GetBooking(bookingID string) (*Booking, error)
	CreateBooking(*Booking) error
	DeleteBooking(bookingID string) error
}
