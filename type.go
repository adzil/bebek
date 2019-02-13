package bebek

import "time"

type Room struct {
	RoomID   string `json:"room_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Booking struct {
	BookingID  string    `json:"booking_id"`
	RoomID     string    `json:"-"`
	Date       time.Time `json:"date"`
	Start      int       `json:"start"`
	End        int       `json:"end"`
	ReservedBy string    `json:"reserved_by"`
}

type Reservation struct {
	*Room
	Booking []*Booking `json:"slots"`
}
