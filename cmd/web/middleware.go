package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/tsawler/vigilate/internal/helpers"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SessionLoad loads the session on requests
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth checks for authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			url := r.URL.Path
			http.Redirect(w, r, fmt.Sprintf("/?target=%s", url), http.StatusFound)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// RecoverPanic recovers from a panic
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Check if there has been a panic
			if err := recover(); err != nil {
				// return a 500 Internal Server response
				helpers.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// NoSurf implements CSRF protection
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.ExemptPath("/pusher/auth")
	csrfHandler.ExemptPath("/pusher/hook")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteStrictMode,
		Domain:   app.Domain,
	})

	return csrfHandler
}

// CheckRemember checks to see if we should log the user in automatically
func CheckRemember(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			cookie, err := r.Cookie(fmt.Sprintf("_%s_gowatcher_remember", preferenceMap["identifier"]))
			if err != nil {
				next.ServeHTTP(w, r)
			} else {
				key := cookie.Value
				// have a remember token, so try to log the user in
				if len(key) > 0 {
					// key length > 0, so it might br a valid token
					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)
					validHash := repo.DB.CheckForToken(id, hash)
					if validHash {
						// valid remember me token, so log the user in
						_ = session.RenewToken(r.Context())
						user, _ := repo.DB.GetUserById(id)
						hashedPassword := user.Password
						session.Put(r.Context(), "userID", id)
						session.Put(r.Context(), "userName", user.FirstName)
						session.Put(r.Context(), "userFirstName", user.FirstName)
						session.Put(r.Context(), "userLastName", user.LastName)
						session.Put(r.Context(), "hashedPassword", string(hashedPassword))
						session.Put(r.Context(), "user", user)
						next.ServeHTTP(w, r)
					} else {
						// invalid token, so delete the cookie
						deleteRememberCookie(w, r)
						session.Put(r.Context(), "error", "You've been logged out from another device!")
						next.ServeHTTP(w, r)
					}
				} else {
					// key length is zero, so it's a leftover cookie (user has not closed browser)
					next.ServeHTTP(w, r)
				}
			}
		} else {
			// they are logged in, but make sure that the remember token has not been revoked
			cookie, err := r.Cookie(fmt.Sprintf("_%s_gowatcher_remember", preferenceMap["identifier"]))
			if err != nil {
				// no cookie
				next.ServeHTTP(w, r)
			} else {
				key := cookie.Value
				// have a remember token, but make sure it's valid
				if len(key) > 0 {
					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)
					validHash := repo.DB.CheckForToken(id, hash)
					if !validHash {
						deleteRememberCookie(w, r)
						session.Put(r.Context(), "error", "You've been logged out from another device!")
						next.ServeHTTP(w, r)
					} else {
						next.ServeHTTP(w, r)
					}
				} else {
					next.ServeHTTP(w, r)
				}
			}
		}
	})
}

// deleteRememberCookie deletes the remember me cookie, and logs the user out
func deleteRememberCookie(w http.ResponseWriter, r *http.Request) {
	_ = session.RenewToken(r.Context())
	// delete the cookie
	newCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_ggowatcher_remember", preferenceMap["identifier"]),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   app.Domain,
		MaxAge:   -1,
		Secure:   app.InProduction,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)

	// log them out
	session.Remove(r.Context(), "userID")
	_ = session.Destroy(r.Context())
	_ = session.RenewToken(r.Context())
}
