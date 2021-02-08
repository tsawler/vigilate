package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/models"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginScreen shows the home (login) screen
func (repo *DBRepo) LoginScreen(w http.ResponseWriter, r *http.Request) {
	// if already logged in, take to dashboard
	if repo.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
		return
	}

	err := helpers.RenderPage(w, r, "login", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Login attempts to log the user in
func (repo *DBRepo) Login(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	id, hash, err := repo.DB.Authenticate(r.Form.Get("email"), r.Form.Get("password"))
	if err == models.ErrInvalidCredentials {
		app.Session.Put(r.Context(), "error", "Invalid login")
		err := helpers.RenderPage(w, r, "login", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
		return
	} else if err == models.ErrInactiveAccount {
		app.Session.Put(r.Context(), "error", "Inactive account!")
		err := helpers.RenderPage(w, r, "login", nil, nil)
		if err != nil {
			printTemplateError(w, err)
		}
		return
	} else if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	if r.Form.Get("remember") == "remember" {
		randomString := helpers.RandomString(12)
		hasher := sha256.New()

		_, err = hasher.Write([]byte(randomString))
		if err != nil {
			log.Println(err)
		}

		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		err = repo.DB.InsertRememberMeToken(id, sha)
		if err != nil {
			log.Println(err)
		}

		// write a cookie
		expire := time.Now().Add(365 * 24 * 60 * 60 * time.Second)
		cookie := http.Cookie{
			Name:     fmt.Sprintf("_%s_gowatcher_remember", app.PreferenceMap["identifier"]),
			Value:    fmt.Sprintf("%d|%s", id, sha),
			Path:     "/",
			Expires:  expire,
			HttpOnly: true,
			Domain:   app.Domain,
			MaxAge:   315360000, // seconds in year
			Secure:   app.InProduction,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
	}

	// we authenticated. Get the user.
	u, err := repo.DB.GetUserById(id)
	if err != nil {
		log.Println(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "userID", id)
	app.Session.Put(r.Context(), "hashedPassword", hash)
	app.Session.Put(r.Context(), "flash", "You've been logged in successfully!")
	app.Session.Put(r.Context(), "user", u)

	if r.Form.Get("target") != "" {
		http.Redirect(w, r, r.Form.Get("target"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
}

// Logout logs the user out
func (repo *DBRepo) Logout(w http.ResponseWriter, r *http.Request) {

	// delete the remember me token, if any
	cookie, err := r.Cookie(fmt.Sprintf("_%s_gowatcher_remember", app.PreferenceMap["identifier"]))
	if err != nil {
	} else {
		key := cookie.Value
		// have a remember token, so get the token
		if len(key) > 0 {
			// key length > 0, so it might be a valid token
			split := strings.Split(key, "|")
			hash := split[1]
			err = repo.DB.DeleteToken(hash)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// delete the remember me cookie, if any
	delCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_gowatcher_remember", app.PreferenceMap["identifier"]),
		Value:    "",
		Domain:   app.Domain,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}
	http.SetCookie(w, &delCookie)

	_ = app.Session.RenewToken(r.Context())
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	repo.App.Session.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
