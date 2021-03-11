package handlers

import (
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
