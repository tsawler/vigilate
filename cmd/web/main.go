package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/pusher/pusher-http-go"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/handlers"
	"github.com/tsawler/vigilate/internal/models"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var app config.AppConfig
var repo *handlers.DBRepo
var session *scs.SessionManager
var preferenceMap map[string]string
var wsClient pusher.Client

const vigilateVersion = "1.0.0"
const maxWorkerPoolSize = 5
const maxJobMaxWorkers = 5

func init() {
	gob.Register(models.User{})
	_ = os.Setenv("TZ", "America/Halifax")
}

// main is the application entry point
func main() {
	// set up application
	insecurePort, err := setupApp()
	if err != nil {
		log.Fatal(err)
	}

	// close channels & db when application ends
	defer close(app.MailQueue)
	defer app.DB.SQL.Close()

	// print info
	log.Printf("******************************************")
	log.Printf("** %sVigilate%s v%s built in %s", "\033[31m", "\033[0m", vigilateVersion, runtime.Version())
	log.Printf("**----------------------------------------")
	log.Printf("** Running with %d Processors", runtime.NumCPU())
	log.Printf("** Running on %s", runtime.GOOS)
	log.Printf("******************************************")

	// create http server
	srv := &http.Server{
		Addr:              *insecurePort,
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Printf("Starting HTTP server on port %s....", *insecurePort)

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
