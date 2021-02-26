package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
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
}

// TestCheck manually tests a host service and sends JSON response
func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")

	log.Println(hostServiceID, oldStatus)

	// get host service
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Service name is", hs.Service.ServiceName)

	// get host?

	// test the service

	// create json
	resp := jsonResp{
		OK:      true,
		Message: "test message",
	}

	// send json to client

	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
