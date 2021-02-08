package handlers

import (
	"github.com/tsawler/vigilate/internal/helpers"
	"net/http"
)

// AllHealthyServices lists all healthy services
func (repo *DBRepo) AllHealthyServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "healthy", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllWarningServices lists all warning services
func (repo *DBRepo) AllWarningServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "warning", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllProblemServices lists all problem services
func (repo *DBRepo) AllProblemServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "problems", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// AllPendingServices lists all pending services
func (repo *DBRepo) AllPendingServices(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "pending", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
