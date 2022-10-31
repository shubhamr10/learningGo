package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/shubhamr10/learningGo/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch expr := mux.(type) {
	case *chi.Mux:
	//test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, type is %T", expr))
	}
}
