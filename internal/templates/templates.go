package templates

import (
	"github.com/tsawler/vigilate/internal/forms"
)

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
