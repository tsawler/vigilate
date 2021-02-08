package templates

import (
	"github.com/tsawler/vigilate/internal/forms"
	"github.com/tsawler/vigilate/internal/models"
)

// TemplateData defines template data
type TemplateData struct {
	CSRFToken       string
	Form            *forms.Form
	IsAuthenticated bool
	PreferenceMap   map[string]string
	User            models.User
	Flash           string
	Warning         string
	Error           string
	GwVersion       string
}
