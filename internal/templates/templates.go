package templates

import (
	"github.com/tsawler/vigilate/internal/forms"
)

// TemplateData defines template data
type TemplateData struct {
	CSRFToken       string
	Form            *forms.Form
	IsAuthenticated bool
	PreferenceMap   map[string]string
	UserID          int
	Flash           string
	Error           string
	GwVersion       string
}
