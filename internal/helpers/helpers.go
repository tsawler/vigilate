package helpers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/models"
	"github.com/tsawler/vigilate/internal/templates"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var app *config.AppConfig
var src = rand.NewSource(time.Now().UnixNano())

// NewHelpers creates new helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// IsAuthenticated returns true if a user is authenticated
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "userID")
	return exists
}

// RandomString returns a random string of letters of length n
func RandomString(n int) string {
	b := make([]byte, n)

	for i, theCache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			theCache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(theCache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		theCache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}

// views is the jet template set
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(),
)

// DefaultData adds default data which is accessible to all templates
func DefaultData(td templates.TemplateData, r *http.Request, w http.ResponseWriter) templates.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = IsAuthenticated(r)
	td.PreferenceMap = app.PreferenceMap
	// if logged in, store user id in template data
	if td.IsAuthenticated {
		u := app.Session.Get(r.Context(), "user").(models.User)
		td.User = u
	}

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	return td
}

// RenderPage renders a page using jet templates
func RenderPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	// add default template data
	var td templates.TemplateData
	if data != nil {
		td = data.(templates.TemplateData)
	}

	// add default data
	td = DefaultData(td, r, w)

	// add template functions
	addTemplateFunctions()

	// load the template and render it
	t, err := views.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
