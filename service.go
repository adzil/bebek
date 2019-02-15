package bebek

import (
	"errors"
	"fmt"
	"time"
)

// GetReservationsRequest represents the GetReservations request data.
type GetReservationsRequest struct {
	RoomID string
	Date   Date
}

// GetSelfReservationsRequest represents the GetSelfReservations request data.
type GetSelfReservationsRequest struct {
	Actor string
}

// CreateBookingRequest represents the CreateBooking request data.
type CreateBookingRequest struct {
	RoomID string `json:"room_id"`
	Date   Date   `json:"date"`
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Actor  string `json:"-"`
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

// Bebek is the default room booking service implementation.
type Bebek struct {
	Repository Repository
}

// GetRooms returns all available room.
func (s *Bebek) GetRooms() ([]*Room, error) {
	return s.Repository.GetRooms()
}

func buildBookingsReservation(reservation *Reservation, bookings []*Booking) {
	for _, booking := range bookings {
		if booking.RoomID == reservation.Room.RoomID {
			reservation.Bookings = append(reservation.Bookings, booking)
		}
	}
}

func buildReservations(rooms []*Room, bookings []*Booking) []*Reservation {
	reservations := make([]*Reservation, len(rooms))
	for i, room := range rooms {
		reservations[i] = &Reservation{
			Room:     room,
			Bookings: make([]*Booking, 0, 8),
		}
		buildBookingsReservation(reservations[i], bookings)
	}
	return reservations
}

// GetReservations returns all reservations with optional filter.
func (s *Bebek) GetReservations(req GetReservationsRequest) ([]*Reservation, error) {
	if req.Date.IsZero() {
		req.Date = Date{time.Now()}
	}
	var rooms []*Room
	var err error
	if req.RoomID != "" {
		var room *Room
		if room, err = s.Repository.GetRoom(req.RoomID); err == nil {
			rooms = []*Room{room}
		}
	} else {
		rooms, err = s.Repository.GetRooms()
	}
	if err != nil {
		return nil, err
	}
	bookings, err := s.Repository.GetBookings(GetBookingsRequest{
		RoomID: req.RoomID,
		Date:   req.Date.Time,
	})
	if err != nil {
		return nil, err
	}
	return buildReservations(rooms, bookings), nil
}

// GetSelfReservations returns all actor's reservations from today and onward.
func (s *Bebek) GetSelfReservations(req GetSelfReservationsRequest) ([]*Reservation, error) {
	rooms, err := s.Repository.GetRooms()
	if err != nil {
		return nil, err
	}
	bookings, err := s.Repository.GetBookings(GetBookingsRequest{
		ReservedBy: req.Actor,
		Date:       time.Now(),
		StartFrom:  true,
	})
	if err != nil {
		return nil, err
	}
	return buildReservations(rooms, bookings), nil
}

// CreateBooking creates new booking that belong to the actor and returns the
// BookingID.
func (s *Bebek) CreateBooking(req CreateBookingRequest) (string, error) {
	startHours := req.Start / 4
	startMins := (req.Start % 4) * 15
	bookingID := req.RoomID + "-" + req.Date.Format("02012006") + "-" + fmt.Sprintf("%02d%02d", startHours, startMins)
	if err := s.Repository.CreateBooking(&Booking{
		BookingID:  bookingID,
		Date:       req.Date,
		Start:      req.Start,
		End:        req.End,
		ReservedBy: req.Actor,
		RoomID:     req.RoomID,
	}); err != nil {
		return "", err
	}
	return bookingID, nil
}

// DeleteBooking deletes BookingID that belongs to the actor. DeleteBooking will
// return error if the booking owner and actor mismatched.
func (s *Bebek) DeleteBooking(req DeleteBookingRequest) error {
	booking, err := s.Repository.GetBooking(req.BookingID)
	if err != nil {
		return err
	}
	if booking.ReservedBy != req.Actor {
		return errors.New("booking not owned by actor")
	}
	return s.Repository.DeleteBooking(req.BookingID)
}
