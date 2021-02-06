package templates

import (
	"github.com/tsawler/vigilate/internal/forms"
	"html/template"
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
