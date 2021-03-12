package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var loginTests = []struct {
	name                 string
	url                  string
	method               string
	postedData           url.Values
	expectedResponseCode int
}{
	{
		name:                 "login-screen",
		url:                  "/",
		method:               "GET",
		expectedResponseCode: http.StatusOK,
	},
	{
		name:   "login-screen-post",
		url:    "/",
		method: "POST",
		postedData: url.Values{
			"email":    {"me@here.com"},
			"password": {"password"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
}

func TestLoginScreen(t *testing.T) {
	for _, e := range loginTests {
		if e.method == "GET" {
			// create a request
			req, _ := http.NewRequest("GET", e.url, nil)

			// add the session info to the context
			ctx := getCtx(req)
			req = req.WithContext(ctx)

			// create a recorder
			rr := httptest.NewRecorder()

			// cast handler we want to test to an http.HandlerFunc
			handler := http.HandlerFunc(Repo.LoginScreen)

			// call the handler with our response recorder (which satisfies the response writer interface),
			// and our request (which has our test session). This executes the method we want to test.
			handler.ServeHTTP(rr, req)

			// check returned status code against expected status code
			if rr.Code != e.expectedResponseCode {
				t.Errorf("%s, expected %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
			}
		} else {
			// create a request with body to post
			req, _ := http.NewRequest("POST", "/", strings.NewReader(e.postedData.Encode()))

			// get our context with the session
			ctx := getCtx(req)
			req = req.WithContext(ctx)

			// create a recorder
			rr := httptest.NewRecorder()

			// cast the handler to a HandlerFunc and call the ServeHTTP method on it.
			// This executes the method we want to test.
			handler := http.HandlerFunc(Repo.Login)
			handler.ServeHTTP(rr, req)

			// test returned status code
			if rr.Code != e.expectedResponseCode {
				t.Errorf("%s, expected %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
			}
		}
	}
}

func TestDBRepo_PusherAuth(t *testing.T) {
	// Now that we have a wrapper for our websocket client, pusher (or ipe)
	// does not have to be running. This is much better.
	postedData := url.Values{
		"socket_id":    {"471281528.421564659"},
		"channel_name": {"private-channel-1"},
	}

	// create the request
	req, _ := http.NewRequest("POST", "/pusher/auth", strings.NewReader(postedData.Encode()))

	// get our context with the session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// create a recorder
	rr := httptest.NewRecorder()

	// cast the handler to a handlerfunc and call serve http
	handler := http.HandlerFunc(Repo.PusherAuth)
	handler.ServeHTTP(rr, req)

	// check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected response 200, but got %d", rr.Code)
	}

	type pusherResp struct {
		Auth string `json:"auth"`
	}

	var p pusherResp

	err := json.NewDecoder(rr.Body).Decode(&p)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Auth) == 0 {
		t.Error("empty json response")
	}
}
