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

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}

// AllReservations returns the list of all the reservations.
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns the list of new the reservations.
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationByID gets reservation by id of the reservation
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}

// UpdateReservations updates a reservation in database
func (m *testDBRepo) UpdateReservations(u models.Reservation) error {

	return nil
}

// DeleteReservation deletes a reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation update a processed reservation
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

// AllRooms fetch all rooms from the database
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}
