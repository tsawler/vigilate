package handlers

import (
	"github.com/pusher/pusher-http-go"
	"github.com/tsawler/vigilate/internal/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// PusherAuth authenticates for pusher
func (repo *DBRepo) PusherAuth(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firstName := app.Session.GetString(r.Context(), "userName")
		userID := app.Session.GetInt(r.Context(), "userID")
		params, _ := ioutil.ReadAll(r.Body)

		presenceData := pusher.MemberData{
			UserID: strconv.Itoa(userID),
			UserInfo: map[string]string{
				"name": firstName,
				"id":   strconv.Itoa(userID),
			},
		}

		response, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
		if err != nil {
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(response)
	}
}

// TestPush test the push functionality
func (repo *DBRepo) TestPush(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		data := make(map[string]string)
		data["message"] = msg
		err := app.WsClient.Trigger("public-channel", "test-event", data)
		if err != nil {
			log.Println(err)
		}
	}
}
