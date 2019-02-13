package bebek

import "time"

// GetReservationsRequest represents the GetReservations request data.
type GetReservationsRequest struct {
	RoomID string
	Date   time.Time
}

// GetSelfReservationsRequest represents the GetSelfReservations request data.
type GetSelfReservationsRequest struct {
	Actor string
}

// CreateBookingRequest represents the CreateBooking request data.
type CreateBookingRequest struct {
	RoomID string
	Date   time.Time
	Start  int
	End    int
	Actor  string
}

// DeleteBookingRequest represents the DeleteBooking request data.
type DeleteBookingRequest struct {
	BookingID string
	Actor     string
}

// Service is the room booking service abstraction.
type Service interface {
	// GetRooms returns all available room.
	GetRooms() ([]*Room, error)

	// GetReservations returns all reservations with optional filter.
	GetReservations(GetReservationsRequest) ([]*Reservation, error)

	// GetSelfReservations returns all actor's reservations from today and
	// onward.
	GetSelfReservations(GetSelfReservationsRequest) ([]*Reservation, error)

	// CreateBooking creates new booking that belong to the actor and returns
	// the BookingID.
	CreateBooking(CreateBookingRequest) (string, error)

	// DeleteBooking deletes BookingID that belongs to the actor. DeleteBooking
	// will return error if the booking owner and actor mismatched.
	DeleteBooking(DeleteBookingRequest) error
}
