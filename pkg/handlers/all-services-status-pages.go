package handlers

import (
	"github.com/tsawler/vigilate/pkg/config"
	"github.com/tsawler/vigilate/pkg/helpers"
	"net/http"
)

// AllHealthyServices lists all healthy services
func (repo *DBRepo) AllHealthyServices(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "healthy", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// AllWarningServices lists all warning services
func (repo *DBRepo) AllWarningServices(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "warning", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// AllProblemServices lists all problem services
func (repo *DBRepo) AllProblemServices(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "problems", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}

// AllPendingServices lists all pending services
func (repo *DBRepo) AllPendingServices(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.RenderPage(w, r, "pending", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
	}
}
