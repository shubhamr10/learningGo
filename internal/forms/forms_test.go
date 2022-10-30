package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/test-url", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when it should have been valid!")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/test-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/test-url", nil)
	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form should show value available")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/test-url", nil)
	form := New(r.PostForm)

	isFirstNameAvailable := form.Has("first_name")
	if isFirstNameAvailable {
		t.Error("first_name should not be available!")
	}

	postedUrl := url.Values{}
	postedUrl.Add("first_name", "Shubham")
	r, _ = http.NewRequest("POST", "/test-url", nil)
	r.PostForm = postedUrl
	form = New(r.PostForm)
	isFirstNameAvailable = form.Has("first_name")
	if !isFirstNameAvailable {
		t.Error("first_name should not be available!")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedUrl := url.Values{}
	form := New(postedUrl)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("email is invalid!")
	}

	postedUrl = url.Values{}
	postedUrl.Add("email", "shubham@acquiew.io")
	form = New(postedUrl)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("email is invalid!")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/test-url", nil)
	form := New(r.PostForm)

	form.MinLength("email", 10)
	if form.Valid() {
		t.Error("email length should be  invalid!")
	}

	postedUrl := url.Values{}
	postedUrl.Add("email", "shubham@acquiew.io")
	r, _ = http.NewRequest("POST", "/test-url", nil)
	r.PostForm = postedUrl
	form = New(r.PostForm)

	form.MinLength("email", 10)
	if !form.Valid() {
		t.Error("email length should be  valid!")
	}
}
