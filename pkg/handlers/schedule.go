package handlers

import (
	"github.com/tsawler/vigilate/pkg/config"
	"github.com/tsawler/vigilate/pkg/helpers"
	"net/http"
)

// ListEntries lists schedule entries
func (repo *DBRepo) ListEntries(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "schedule", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}
