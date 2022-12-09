package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shubhamr10/learningGo/internal/config"
	"github.com/shubhamr10/learningGo/internal/driver"
	"github.com/shubhamr10/learningGo/internal/forms"
	"github.com/shubhamr10/learningGo/internal/models"
	"github.com/shubhamr10/learningGo/internal/render"
	"github.com/shubhamr10/learningGo/internal/repository"
	"github.com/shubhamr10/learningGo/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo is the variable used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository for tests
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Generals is the homepage handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the homepage handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// MakeReservations is the homepage handler
func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		//helpers.ServerError(w, errors.New("cannot get reservations from session"))
		//return
	}
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		//helpers.ServerError(w, err)
		//return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostMakeReservations accepts the post request of reservation form
func (m *Repository) PostMakeReservations(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, errors.New("cannot get from session"))
		// return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)
	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//Insert data to reservation table
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert into database.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room roomaretiction.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}

	// send email notification - first to guest
	htmlMessage := fmt.Sprintf(`
			<strong>Reservation Confirmation</strong><br>
			Dear %s, <br>
			This is to confirm that your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	// send email notification to - property owner
	htmlMessage = fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear Shubham, <br>
		This is to confirm that Mr./Mrs. %s has a booking for room no. %d from  %s to %s.
`, reservation.FirstName, reservation.RoomID, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg = models.MailData{
		To:      "property-owner@acq.io",
		From:    "me@here.com",
		Subject: "Reservation confirmation",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// SearchAvailability is the homepage handler
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability is the homepage handler
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	// convert into time.Layout format
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid start date.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid end date.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error while fetching availability.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability!")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// PostSearchAvailabilityJSON handles request for availability and send JSON response
func (m *Repository) PostSearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return apprpriate json
		resp := jsonResponse{
			OK:      false,
			Message: "internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	available, err4 := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err4 != nil {
		// can't parse form, so return apprpriate json
		resp := jsonResponse{
			OK:      false,
			Message: "internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "	")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "!",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}
	out, _ := json.MarshalIndent(resp, "", "     ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the homepage handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the reservations summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//log.Println("cannot get item from session")
		// m.App.ErrorLog.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// used to have next 6 lines
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	log.Println(err)
	//	m.App.Session.Put(r.Context(), "error", "missing url parameter")
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}

	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "unable to find reservation to put room upon")
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)
}

// BookRoom takes url parameters, builds a session variable and takes user to make-reservations screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter for room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err)
		// return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	room, err3 := m.DB.GetRoomByID(roomId)
	if err3 != nil {
		m.App.Session.Put(r.Context(), "error", "invalid room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
		// helpers.ServerError(w, err3)
		// return
	}
	var res models.Reservation
	res.RoomID = roomId
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)

}

// ShowLogin is used to display the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin is used to receive the login request
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println("err", err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println("err", err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout is used to log out a user
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
