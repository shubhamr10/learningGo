package dbrepo

import (
	"errors"
	"time"

	"github.com/shubhamr10/learningGo/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return false
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2 then fail, otherwise pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID return true if rooms are available, and false if no rooms are available (for particular roomID)
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	layout := "2006-01-02"
	booked_date_start, _ := time.Parse(layout, "2050-01-01")
	booked_date_end, _ := time.Parse(layout, "2050-02-01")
	// if room is invalid thow an error
	if roomID == 100 {
		return false, errors.New("some error")
	}
	// if the start and end lies then no room is available
	if start.After(booked_date_start) && start.Before(booked_date_end) {
		// no rooms available
		return false, nil
	}
	// else available
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of rooms for any date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	layout := "2006-01-02"
	booked_date_start, _ := time.Parse(layout, "2050-01-01")
	booked_date_end, _ := time.Parse(layout, "2050-02-01")
	// if the start and end lies then no room is available
	if start.After(booked_date_start) && start.Before(booked_date_end) {
		// no rooms available
		return rooms, nil
	}
	if start == booked_date_start {
		return rooms, errors.New("some error")
	}
	// roms are available
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}
