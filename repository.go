package bebek

import (
	"database/sql"
	"errors"
	"time"
)

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

	// GetRoom return Room by their ID
	GetRoom(roomID string) (*Room, error)

	// GetBookings returns all booking that matches the request filter.
	GetBookings(GetBookingsRequest) ([]*Booking, error)

	// GetBooking returns Booking from its BookingID
	GetBooking(bookingID string) (*Booking, error)

	// CreateBooking creates a new Booking.
	CreateBooking(*Booking) error

	// DeleteBooking deletes a booking from repository.
	DeleteBooking(bookingID string) error
}

type MySQLRepository struct {
	DB *sql.DB
}

func (r *MySQLRepository) GetRooms() ([]*Room, error) {
	rows, err := r.DB.Query("SELECT `a`.`room_id`, `a`.`name`, `a`.`location` FROM `room` AS `a`")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]*Room, 0)
	for rows.Next() {
		room := &Room{}
		if err := rows.Scan(&room.RoomID, &room.Name, &room.Location); err != nil {
			return nil, err
		}

		results = append(results, room)
	}

	return results, nil
}

func (r *MySQLRepository) GetRoom(roomID string) (*Room, error) {
	rows, err := r.DB.Query("SELECT `a`.`room_id`, `a`.`name`, `a`.`location` FROM room AS a WHERE `a`.`room_id` = ?", roomID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	room := &Room{}
	if !rows.Next() {
		return nil, errors.New("Unable to find Room with ID: " + roomID)
	}

	if err = rows.Scan(&room.RoomID, &room.Name, &room.Location); err != nil {
		return nil, err
	}

	return room, nil
}

func (r *MySQLRepository) GetBookings(req GetBookingsRequest) ([]*Booking, error) {
	query := "SELECT `a`.`booking_id`, `a`.`room_id`, `a`.`date`, `a`.`slot`, `a`.`reserved_by` FROM booking AS a "
	args := []interface{}{req.Date}
	if req.StartFrom {
		query += "WHERE `a`.`date` >= DATE(?) "
	} else {
		query += "WHERE `a`.`date` = DATE(?) "
	}
	if req.ReservedBy != "" {
		query += "AND `a`.`reserved_by` = ? "
		args = append(args, req.ReservedBy)
	}
	if req.RoomID != "" {
		query += "AND a.room_id = ? "
		args = append(args, req.RoomID)
	}
	query += "ORDER BY a.room_id, a.date, a.slot"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []*Booking

	var previousBookingID string
	var booking *Booking
	for rows.Next() {
		var (
			bookingID  string
			roomID     string
			date       time.Time
			slot       int
			reservedBy string
		)

		if err := rows.Scan(&bookingID, &roomID, &date, &slot, &reservedBy); err != nil {
			return nil, err
		}

		if booking == nil || previousBookingID != bookingID {
			booking = &Booking{
				BookingID:  bookingID,
				RoomID:     roomID,
				Date:       Date{date},
				Start:      slot,
				ReservedBy: reservedBy,
			}

			results = append(results, booking)
			previousBookingID = bookingID
		}

		booking.End = slot
	}

	return results, nil
}

func (r *MySQLRepository) GetBooking(bookingID string) (*Booking, error) {
	rows, err := r.DB.Query("SELECT `a`.`room_id`, `a`.`date`, `a`.`slot`, `a`.`reserved_by` FROM booking AS a WHERE `a`.`booking_id` = ? ORDER BY `a`.`slot`", bookingID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var booking *Booking
	for rows.Next() {
		var (
			roomID     string
			date       time.Time
			slot       int
			reservedBy string
		)

		if err := rows.Scan(&roomID, &date, &slot, &reservedBy); err != nil {
			return nil, err
		}

		if booking == nil {
			booking = &Booking{
				BookingID:  bookingID,
				RoomID:     roomID,
				Date:       Date{date},
				Start:      slot,
				ReservedBy: reservedBy,
			}
		}

		booking.End = slot
	}

	return booking, nil
}

func (r *MySQLRepository) CreateBooking(booking *Booking) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	for i := booking.Start; i <= booking.End; i++ {
		_, err := tx.Exec("INSERT INTO booking(`room_id`, `date`, `slot`, `booking_id`, `reserved_by`) VALUES(?, ?, ?, ?, ?)", booking.RoomID, booking.Date.Time, i, booking.BookingID, booking.ReservedBy)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *MySQLRepository) DeleteBooking(bookingID string) error {
	_, err := r.DB.Exec("DELETE FROM booking WHERE `booking_id` = ?", bookingID)

	return err
}
