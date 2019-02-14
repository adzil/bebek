package bebek

import "time"

// Room represents the room information.
type Room struct {
	RoomID   string `json:"room_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Booking represents the booking information.
type Booking struct {
	BookingID  string    `json:"booking_id"`
	RoomID     string    `json:"-"`
	Date       time.Time `json:"date"`
	Start      int       `json:"start"`
	End        int       `json:"end"`
	ReservedBy string    `json:"reserved_by"`
}

// Reservation represents the room and booking information.
type Reservation struct {
	*Room
	Bookings []*Booking `json:"bookings"`
}
