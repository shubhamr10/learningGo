package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/shubhamr10/learningGo/internal/models"
)

// AppConfig Holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
