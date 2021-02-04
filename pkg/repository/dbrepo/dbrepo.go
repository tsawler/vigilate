package dbrepo

import (
	"database/sql"
	"github.com/tsawler/vigilate/pkg/config"
	"github.com/tsawler/vigilate/pkg/repository"
)

var app *config.AppConfig

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates the repository
func NewPostgresRepo(Conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	app = a
	return &postgresDBRepo{
		App: a,
		DB:  Conn,
	}
}
