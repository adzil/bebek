package bebek

import (
	"fmt"
	"strings"
	"time"
)

// DateLayout is time format used for parsing date string
const DateLayout = "020106"

var nilTime = (time.Time{}).UnixNano()

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
	if d.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format(DateLayout))), nil
}

// IsSet check whether is value of Date has been set
func (d *Date) IsSet() bool {
	return d.UnixNano() != nilTime
}
