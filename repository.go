package bebek

import "time"

// GetBookingsRequest represents the GetBookings request data.
type GetBookingsRequest struct {
	RoomID     string
	Date       time.Time
	StartFrom  bool
	ReservedBy string
}

// Repository is the room booking repository abstraction.
type Repository interface {
	// GetRooms returns all available room.
	GetRooms() ([]*Room, error)

	// GetBookings returns all booking that matches the request filter.
	GetBookings(GetBookingsRequest) ([]*Booking, error)

	// GetBooking returns Booking from its BookingID
	GetBooking(bookingID string) (*Booking, error)

	// CreateBooking creates a new Booking.
	CreateBooking(*Booking) error

	// DeleteBooking deletes a booking from repository.
	DeleteBooking(bookingID string) error
}
