package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shubhamr10/learningGo/internal/driver"
	"github.com/shubhamr10/learningGo/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarter", "/generals-quarters", "GET", http.StatusOK},
	{"major's-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"make-reservations", "/make-reservations", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//{"post-search-availability", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2020-02-03"},
	//	{key: "end", value: "2020-02-08"},
	//}, http.StatusOK},
	//{"post-search-availability", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "2020-02-03"},
	//	{key: "end", value: "2020-02-08"},
	//}, http.StatusOK},
	//{"make-reservation-post", "/make-reservations", "POST", []postData{
	//	{key: "first_name", value: "John"},
	//	{key: "last_name", value: "Smith"},
	//	{key: "email", value: "m2@here.com"},
	//	{key: "phone", value: "555-555-5555"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	// table tests
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
			return
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_MakeReservations(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarter",
		},
	}

	req, _ := http.NewRequest(http.MethodGet, "/make-reservations", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	// rr stands for responseRecorder
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest(http.MethodGet, "/make-reservation", nil)
	ctx = GetCtx(req)

	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest(http.MethodGet, "/make-reservation", nil)
	ctx = GetCtx(req)

	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostMakeReservations(t *testing.T) {
	// # test case 1 - where data is present
	// convert into time.Layout format
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2050-01-01")
	endDate, _ := time.Parse(layout, "2050-02-01")

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	reqBody := "first_name=John"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Make_Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Make_Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid date and room id was there in the video
	// but we are using session which are pre-defined then we skip those tests

	// test for invalid session
	// write table test instead
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Make_Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// running test for invalid form data
	reqBody = "first_name=Jn"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=johnsmith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post Make_Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test insert insert reservation
	reservation.RoomID = 2
	reqBody = "first_name=John"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Make_Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure of room restriction entry to database
	reservation.RoomID = 1000
	reqBody = "first_name=John"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostMakeReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Make_Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostSearchAvailabilityJSON(t *testing.T) {
	// # test cases
	/*
		* test cases : rooms are not available
		test case 2: rooms are available
	*/
	// test case : rooms are not available
	reqBody := "start=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest(http.MethodPost, "/search-availability-json", strings.NewReader(reqBody))
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var jsonRes jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &jsonRes)
	if err != nil {
		t.Error("failed to parse jsonRes json")
	}

	if jsonRes.OK {
		t.Error("date should not be available but got availability")
	}
	// test case : roooms are available
	reqBody = "start=2022-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest(http.MethodPost, "/search-availability-json", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonRes)
	if err != nil {
		t.Error("failed to parse jsonRes json")
	}

	if !jsonRes.OK {
		t.Error("date should be available but got no availability")
	}

	// test case : roooms are not available and generate an error
	reqBody = "start=2022-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=100")

	req, _ = http.NewRequest(http.MethodPost, "/search-availability-json", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonRes)
	if err != nil {
		t.Error("failed to parse jsonRes json")
	}

	if jsonRes.OK {
		t.Error("date should be not available but got availability")
	}

	// test case 4 : no form body provided
	req, _ = http.NewRequest(http.MethodPost, "/search-availability-json", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostSearchAvailabilityJSON)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &jsonRes)
	if err != nil {
		t.Error("failed to parse jsonRes json")
	}

	if jsonRes.OK {
		t.Error("date should be not available but got availability")
	}
}

func TestRepository_PostSearchAvailability(t *testing.T) {
	// test case : form is invalid
	req, _ := http.NewRequest(http.MethodPost, "/make-reservations", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostSearchAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case: start date is invalid
	reqBody := "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-05")
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostSearchAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case: end date is invalid
	reqBody = "start=2050-01-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostSearchAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case : rooms not available
	reqBody = "start=2050-01-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-02-03")
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostSearchAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case: rooms are available
	reqBody = "start=2030-01-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2030-02-03")
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostSearchAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case : rooms not available with an error
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-02-03")
	req, _ = http.NewRequest(http.MethodPost, "/make-reservations", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	// setting headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// rr stands for responseRecorder
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostSearchAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostSearchAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_NewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	// test case : No session data, should return an error
	req, _ := http.NewRequest(http.MethodGet, "/reservation-summary", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	// I need a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation Summary returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case : Session is present, should work as expected
	layout := "2006-01-02"
	booking_date_start, _ := time.Parse(layout, "2050-01-01")
	booking_date_end, _ := time.Parse(layout, "2050-02-01")
	reservation := models.Reservation{
		StartDate: booking_date_start,
		EndDate:   booking_date_end,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@smith.com",
		Phone:     "1234567789",
		RoomID:    1,
	}
	req, _ = http.NewRequest(http.MethodGet, "reservation-summary", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation Summary returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// reservation := models.Reservation{
	// 	RoomID: 1,
	// 	Room: models.Room{
	// 		ID:       1,
	// 		RoomName: "General's Quarter",
	// 	},
	// }

	// req, _ := http.NewRequest(http.MethodGet, "/make-reservations", nil)
	// ctx := GetCtx(req)
	// req = req.WithContext(ctx)

	// // rr stands for responseRecorder
	// rr := httptest.NewRecorder()
	// session.Put(ctx, "reservation", reservation)

	// handler := http.HandlerFunc(Repo.MakeReservations)
	// handler.ServeHTTP(rr, req)

	// if rr.Code != http.StatusOK {
	// 	t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	// }
}

func TestRepository_ChooseRoom(t *testing.T) {
	// test case : valid data
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case: invalid room
	req, _ = http.NewRequest(http.MethodGet, "/choose-room/", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/"

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("choose room handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case : session for reservation not present in the code
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	// test case : everythigng is valid
	req, _ := http.NewRequest(http.MethodGet, "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("book rooom handler returned a wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case : database error for room id
	req, _ = http.NewRequest(http.MethodGet, "/book-room?s=2050-01-01&e=2050-01-02&id=100", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("book rooom handler returned a wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case room id is invalid
	req, _ = http.NewRequest(http.MethodGet, "/book-room?s=2050-01-01&e=2050-01-02&id=room", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("book rooom handler returned a wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func GetCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
