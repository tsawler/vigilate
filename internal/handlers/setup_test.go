package handlers

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
	"github.com/tsawler/vigilate/internal/channeldata"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/driver"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/repository/dbrepo"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var testSession *scs.SessionManager

func TestMain(m *testing.M) {

	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false

	mailQueue := make(chan channeldata.MailJob, 5)

	// define application configuration
	a := config.AppConfig{
		DB:           &driver.DB{},
		Session:      testSession,
		InProduction: false,
		Domain:       "localhost",
		MailQueue:    mailQueue,
	}

	app = &a

	preferenceMap := make(map[string]string)

	app.PreferenceMap = preferenceMap

	// create pusher client
	dws := dummyWs{
		AppID:  "1",
		Secret: "123abc",
		Key:    "abc123",
		Secure: false,
		Host:   "localhost:4001",
	}

	app.WsClient = &dws

	monitorMap := make(map[int]cron.EntryID)
	app.MonitorMap = monitorMap

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	app.Scheduler = scheduler

	repo := NewTestHandlers(app)
	NewHandlers(repo, app)

	helpers.NewHelpers(app)

	helpers.SetViews("./../../views")

	os.Exit(m.Run())
}

// gets the context with session added
func getCtx(req *http.Request) context.Context {
	ctx, err := testSession.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

// NewTestHandlers creates a new repository
func NewTestHandlers(a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// dummyWs is a type that satisfies the pusher.Client interface
type dummyWs struct {
	AppID                        string
	Key                          string
	Secret                       string
	Host                         string // host or host:port pair
	Secure                       bool   // true for HTTPS
	Cluster                      string
	HTTPClient                   *http.Client
	EncryptionMasterKey          string  // deprecated
	EncryptionMasterKeyBase64    string  // for E2E
	validatedEncryptionMasterKey *[]byte // parsed key for use
}

func (c *dummyWs) Trigger(channel string, eventName string, data interface{}) error {
	return nil
}

func (c *dummyWs) TriggerMulti(channels []string, eventName string, data interface{}) error {
	return nil
}

func (c *dummyWs) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) error {
	return nil
}

func (c *dummyWs) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) error {
	return nil
}

func (c *dummyWs) TriggerBatch(batch []pusher.Event) error {
	return nil
}

func (c *dummyWs) Channels(additionalQueries map[string]string) (*pusher.ChannelsList, error) {
	var cl pusher.ChannelsList
	return &cl, nil
}

func (c *dummyWs) Channel(name string, additionalQueries map[string]string) (*pusher.Channel, error) {
	var cl pusher.Channel
	return &cl, nil
}

func (c *dummyWs) GetChannelUsers(name string) (*pusher.Users, error) {
	var cl pusher.Users
	return &cl, nil
}

func (c *dummyWs) AuthenticatePrivateChannel(params []byte) (response []byte, err error) {
	return []byte("Hello"), nil
}

func (c *dummyWs) AuthenticatePresenceChannel(params []byte, member pusher.MemberData) (response []byte, err error) {
	jStr := `
	{"auth":"abc123:b75c9f83f1a1dbe7d6933316348039c6270b27a416286385a0fed98529cf46d1","channel_data":"{\"user_id\":\"1\",\"user_info\":{\"id\":\"1\",\"name\":\"Admin\"}}"}
`
	return []byte(jStr), nil
}

func (c *dummyWs) Webhook(header http.Header, body []byte) (*pusher.Webhook, error) {
	var wh pusher.Webhook
	return &wh, nil
}
