package bebek

import (
	"fmt"
	"strings"
	"time"
)

// DateLayout is time format used for parsing date string
const DateLayout = "020106"

// Date is time.Time which layout using DateLayout
type Date struct {
	time.Time
}

// Room represents the room information.
type Room struct {
	RoomID   string `json:"room_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Booking represents the booking information.
type Booking struct {
	BookingID  string `json:"booking_id"`
	RoomID     string `json:"-"`
	Date       Date   `json:"date"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	ReservedBy string `json:"reserved_by"`
}

// Reservation represents the room and booking information.
type Reservation struct {
	*Room
	Bookings []*Booking `json:"bookings"`
}

// UnmarshalText decode date string into date.
func (d *Date) UnmarshalText(b []byte) error {
	var err error
	d.Time, err = time.Parse(DateLayout, string(b))
	return err
}

// MarshalText encode date into string.
func (d *Date) MarshalText() ([]byte, error) {
	if d.IsZero() {
		return []byte{}, nil
	}
	return []byte(d.Time.Format(DateLayout)), nil
}

// UnmarshalJSON decode Date as JSON
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return
	}
	d.Time, err = time.Parse(DateLayout, s)
	return
}

// MarshalJSON encode Date to JSON
func (d *Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format(DateLayout))), nil
}
