package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/robfig/cron/v3"
	"github.com/tsawler/vigilate/internal/channeldata"
	"github.com/tsawler/vigilate/internal/driver"
	"github.com/tsawler/vigilate/internal/models"
	"html/template"
)

// AppConfig holds application configuration
type AppConfig struct {
	DB            *driver.DB
	Session       *scs.SessionManager
	InProduction  bool
	Domain        string
	MonitorMap    map[int]cron.EntryID
	PreferenceMap map[string]string
	Scheduler     *cron.Cron
	WsClient      models.WSClient
	PusherSecret  string
	TemplateCache map[string]*template.Template
	MailQueue     chan channeldata.MailJob
	Version       string
	Identifier    string
}
