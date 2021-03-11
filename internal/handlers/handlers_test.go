package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDBRepo_LoginScreen(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AdminDashboard)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed login screen: expected code 200, but got %d", rr.Code)
	}
}

// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
