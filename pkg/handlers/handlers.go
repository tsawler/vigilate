package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/tsawler/vigilate/pkg/config"
	"github.com/tsawler/vigilate/pkg/driver"
	"github.com/tsawler/vigilate/pkg/helpers"
	"github.com/tsawler/vigilate/pkg/repository"
	"github.com/tsawler/vigilate/pkg/repository/dbrepo"
	"log"
	"net/http"
	"runtime/debug"
)

//Repo is the repository
var Repo *DBRepo
var app *config.AppConfig

// DBRepo is the db repo
type DBRepo struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewHandlers creates the handlers
func NewHandlers(repo *DBRepo, a *config.AppConfig) {
	Repo = repo
	app = a
}

// NewPostgresqlHandlers creates db repo for postgres
func NewPostgresqlHandlers(db *driver.DB, a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// AdminDashboard displays the dashboard
func (repo *DBRepo) AdminDashboard(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := make(jet.VarMap)
		vars.Set("no_healthy", 0)
		vars.Set("no_problem", 0)
		vars.Set("no_pending", 0)
		vars.Set("no_warning", 0)

		err := helpers.RenderPage(w, r, "dashboard", vars, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// Events displays the events page
func (repo *DBRepo) Events(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "events", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// Settings displays the settings page
func (repo *DBRepo) Settings(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "settings", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// PostSettings saves site settings
func (repo *DBRepo) PostSettings(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefMap := make(map[string]string)

		prefMap["site_url"] = r.Form.Get("site_url")
		prefMap["notify_name"] = r.Form.Get("notify_name")
		prefMap["notify_email"] = r.Form.Get("notify_email")
		prefMap["smtp_server"] = r.Form.Get("smtp_server")
		prefMap["smtp_port"] = r.Form.Get("smtp_port")
		prefMap["smtp_user"] = r.Form.Get("smtp_user")
		prefMap["smtp_password"] = r.Form.Get("smtp_password")
		prefMap["sms_enabled"] = r.Form.Get("sms_enabled")
		prefMap["sms_provider"] = r.Form.Get("sms_provider")
		prefMap["twilio_phone_number"] = r.Form.Get("twilio_phone_number")
		prefMap["twilio_sid"] = r.Form.Get("twilio_sid")
		prefMap["twilio_auth_token"] = r.Form.Get("twilio_auth_token")
		prefMap["smtp_from_email"] = r.Form.Get("smtp_from_email")
		prefMap["smtp_from_name"] = r.Form.Get("smtp_from_name")
		prefMap["notify_via_sms"] = r.Form.Get("notify_via_sms")
		prefMap["notify_via_email"] = r.Form.Get("notify_via_email")
		prefMap["sms_notify_number"] = r.Form.Get("sms_notify_number")

		if r.Form.Get("sms_enabled") == "0" {
			prefMap["notify_via_sms"] = "0"
		}

		err := repo.DB.InsertOrUpdateSitePreferences(prefMap)
		if err != nil {
			log.Println(err)
			ClientError(w, r, http.StatusBadRequest)
			return
		}

		// update app config
		for k, v := range prefMap {
			app.PreferenceMap[k] = v
		}

		app.Session.Put(r.Context(), "flash", "Changes saved")

		if r.Form.Get("action") == "1" {
			http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
		}
	}
}

// AllHosts displays list of all hosts
func (repo *DBRepo) AllHosts(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "hosts", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// Host shows the host add/edit form
func (repo *DBRepo) Host(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "host", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// ClientError will display error page for client error i.e. bad request
func ClientError(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		show404(w, r)
	case http.StatusInternalServerError:
		show500(w, r)
	default:
		http.Error(w, http.StatusText(status), status)
	}
}

// Show404 shows a 404 page
func (repo *DBRepo) Show404(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		show404(w, r)
	}
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)
	show500(w, r)
}

func show404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/404.html")
	return
}

func show500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}

func printTemplateError(w http.ResponseWriter, err error) {
	_, _ = fmt.Fprint(w, fmt.Sprintf(`<small><span class='text-danger'>Error executing template: %s</span></small>`, err))
}
