package templates

import (
	"fmt"
	"github.com/tsawler/vigilate/pkg/config"
	"github.com/tsawler/vigilate/pkg/forms"
	"html/template"
	"path/filepath"
	"strconv"
	"time"
)

var templatePath string

var functions = template.FuncMap{
	"formatDateWithLayout": FormatDateWithLayout,
	"dateAfterYearOne":     DateAfterY1,
	"thisYear":             ThisYear,
}

// TemplateData defines template data
type TemplateData struct {
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	RowSets         map[string]interface{}
	CSRFToken       string
	Form            *forms.Form
	IsAuthenticated bool
	PreferenceMap   map[string]string
	UserID          int
	Flash           string
	Error           string
	GwVersion       string
}

// NewTemplateCache creates the template cache
func NewTemplateCache(app *config.AppConfig) (map[string]*template.Template, error) {
	templatePath = "./ui/html"
	mailTemplatePath := "./ui/mail"
	myCache := map[string]*template.Template{}

	// pages
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", templatePath))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", templatePath))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", templatePath))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}

		matches, err = filepath.Glob(fmt.Sprintf("%s/*.partial.tmpl", templatePath))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.partial.tmpl", templatePath))
			if err != nil {
				return nil, err
			}
		}

		matches, err = filepath.Glob(fmt.Sprintf("%s/partials/*.partial.tmpl", templatePath))
		if err != nil {
			return nil, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/partials/*.partial.tmpl", templatePath))
			if err != nil {
				return nil, err
			}
		}

		// Add the template set to the cache,
		myCache[name] = ts

		//myCache["top-bar"] = ts

		// now do mail templates
		mails, err := filepath.Glob(fmt.Sprintf("%s/*.mail.tmpl", mailTemplatePath))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, page := range mails {
			name := filepath.Base(page)

			ts, err := template.New(name).Funcs(functions).ParseFiles(page)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", mailTemplatePath))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if len(matches) > 0 {
				ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", mailTemplatePath))
				if err != nil {
					return nil, err
				}
			}

			matches, err = filepath.Glob(fmt.Sprintf("%s/*.partial.tmpl", mailTemplatePath))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if len(matches) > 0 {
				ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.partial.tmpl", mailTemplatePath))
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
			}

			// Add the template set to the cache,
			myCache[name] = ts

		}
	}
	app.TemplateCache = myCache
	return myCache, nil
}

// FormatDateWithLayout formats a date/time with specified layout string
func FormatDateWithLayout(t time.Time, f string) string {
	return t.Format(f)
}

// DateAfterY1 is used to verify that a date is after the year 1 (since go hates nulls)
func DateAfterY1(t time.Time) bool {
	yearOne := time.Date(0001, 11, 17, 20, 34, 58, 651387237, time.UTC)
	return t.After(yearOne)
}

// ThisYear returns current year as YYYY
func ThisYear() string {
	return strconv.Itoa(time.Now().Year())
}
