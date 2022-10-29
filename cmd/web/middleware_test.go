package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch expr := h.(type) {
	case http.Handler:
	//do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", expr))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch expr := h.(type) {
	case http.Handler:
	//do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", expr))
	}
}
