package dbrepo

import (
	"context"
	"github.com/shubhamr10/learningGo/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return false
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	stmt := "INSERT INTO reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6,$7, $8, $9) returning id"
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restrictions_id) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID return true if rooms are available, and false if no rooms are available (for particular roomID)
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT COUNT(id) FROM room_restrictions " +
		"WHERE room_id = $1 and start_date < $2 and end_date > $3 ;"
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, roomID, end, start)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, err
}

// SearchAvailabilityForAllRooms returns a slice of rooms for any date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := "SELECT r.id, r.room_name from rooms r where r.id not in" +
		" (SELECT room_id from room_restrictions rr where rr.start_date < $1 and end_date > $2)"

	rows, err := m.DB.QueryContext(ctx, query, end, start)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := "SELECT id, room_name, created_at, updated_at from rooms where id = $1"
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}
