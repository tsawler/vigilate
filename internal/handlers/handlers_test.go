package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestDBRepo_LoginScreen tests showing the login screen
func TestDBRepo_LoginScreen(t *testing.T) {
	// create a request
	req, _ := http.NewRequest("GET", "/", nil)

	// add the session info to the context
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// create a recorder
	rr := httptest.NewRecorder()

	// cast handler we want to test to an http.HandlerFunc
	handler := http.HandlerFunc(Repo.LoginScreen)

	// call the handler with our response recorder (which satisfies the response writer interface),
	// and our request (which has our test session)
	handler.ServeHTTP(rr, req)

	// check returned status code against expectd status code
	if rr.Code != http.StatusOK {
		t.Errorf("failed login screen: expected 200, but got %s - %d", rr.Code)
	}
}

// TestDBRepo_Login tests logging in
func TestDBRepo_Login(t *testing.T) {
	// put create url.Values with posted data
	postedData := url.Values{
		"email":    {"admin@example.com"},
		"password": {"password"},
	}

	// create a request with body to post
	req, _ := http.NewRequest("POST", "/", strings.NewReader(postedData.Encode()))

	// get our context with the session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// create a recorder
	rr := httptest.NewRecorder()

	// cast the handler to a handlerfunc and call serve http
	handler := http.HandlerFunc(Repo.Login)
	handler.ServeHTTP(rr, req)

	// test returned status code
	if rr.Code != http.StatusSeeOther {
		t.Errorf("failed post login screen: expected 303, but got %d", rr.Code)
	}
}
