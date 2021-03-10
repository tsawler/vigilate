package handlers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
	"github.com/tsawler/vigilate/internal/channeldata"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/driver"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/repository/dbrepo"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager

func TestMain(m *testing.M) {
	// session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	mailQueue := make(chan channeldata.MailJob, 5)

	// define application configuration
	a := config.AppConfig{
		DB:           &driver.DB{},
		Session:      session,
		InProduction: false,
		Domain:       "localhost",
		MailQueue:    mailQueue,
	}

	repo := NewTestHandlers(app)
	NewHandlers(repo, app)

	app = &a

	preferenceMap := make(map[string]string)

	app.PreferenceMap = preferenceMap

	// create pusher client
	wsClient := pusher.Client{
		AppID:  "1",
		Secret: "123abc",
		Key:    "abc123",
		Secure: false,
		Host:   "localhost:4001",
	}

	app.WsClient = wsClient

	monitorMap := make(map[int]cron.EntryID)
	app.MonitorMap = monitorMap

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	app.Scheduler = scheduler

	helpers.NewHelpers(app)

	os.Exit(m.Run())
}

// NewTestHandlers creates a new repository
func NewTestHandlers(a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}
