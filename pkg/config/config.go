package config

import (
	"html/template"
	"log"
)

// AppConfig Holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
}
