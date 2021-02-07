package main

import (
	"github.com/go-chi/chi"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/handlers"
	"net/http"
)

func routes(app config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	// default middleware
	mux.Use(SessionLoad)
	mux.Use(RecoverPanic)
	mux.Use(NoSurf)
	mux.Use(CheckRemember)

	// login
	mux.Get("/", handlers.Repo.LoginScreen(app))
	mux.Post("/", handlers.Repo.Login(app))

	mux.Get("/user/logout", handlers.Repo.Logout(app))

	mux.Get("/ws/test-push", handlers.Repo.TestPush(app))

	// pusher routes - excluded from csrf protection (nosurf)
	mux.Route("/pusher", func(mux chi.Router) {
		// pusher route requires authentication
		mux.Use(Auth)
		mux.Post("/auth", handlers.Repo.PusherAuth(app))
	})

	// admin routes
	mux.Route("/admin", func(mux chi.Router) {
		// all admin routes are protected
		mux.Use(Auth)

		// overview
		mux.Get("/overview", handlers.Repo.AdminDashboard(app))

		// events
		mux.Get("/events", handlers.Repo.Events(app))

		// settings
		mux.Get("/settings", handlers.Repo.Settings(app))
		mux.Post("/settings", handlers.Repo.PostSettings(app))

		// service status pages (all hosts)
		mux.Get("/all-healthy", handlers.Repo.AllHealthyServices(app))
		mux.Get("/all-warning", handlers.Repo.AllWarningServices(app))
		mux.Get("/all-problems", handlers.Repo.AllProblemServices(app))
		mux.Get("/all-pending", handlers.Repo.AllPendingServices(app))

		// schedule
		mux.Get("/schedule", handlers.Repo.ListEntries(app))

		// hosts
		mux.Get("/host/all", handlers.Repo.AllHosts(app))
		mux.Get("/host/{id}", handlers.Repo.Host(app))

		mux.Handle("/*", handlers.Repo.Show404(app))
	})

	// static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// everything else is 404
	mux.Handle("/*", handlers.Repo.Show404(app))

	return mux
}
