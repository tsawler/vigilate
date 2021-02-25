package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/tsawler/vigilate/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// jsonResp describes the JSON response sent back to client
type jsonResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
	IsError       bool      `json:"is_error"`
}

// TestCheck manually tests a host service and sends JSON response
func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")

	// get host service
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	// get host?
	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	// test the service
	msg, newStatus, _ := repo.testServiceForHost(hs, h)

	isError := false

	if newStatus != "healthy" {
		isError = true
	}
	// create json
	resp := jsonResp{
		OK:        true,
		Message:   msg,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		LastCheck: time.Now(),
		IsError:   isError,
	}

	// TODO: if old status != new status, then update the database

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// testServiceForHost tests a service for a host
func (repo *DBRepo) testServiceForHost(hs models.HostService, h models.Host) (string, string, error) {
	var msg, newStatus string

	switch hs.Service.ServiceName {
	case "HTTP":
		msg, newStatus = repo.testHTTPForHostService(h.URL)
		break

	default:
	}

	return msg, newStatus, nil
}

// testHTTPForHostService tests http service for a host
func (repo *DBRepo) testHTTPForHostService(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	// make sure this is an http request
	url = strings.Replace(url, "https://", "http://", -1)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}
