package dbrepo

import (
	"context"
	"github.com/tsawler/vigilate/internal/models"
	"log"
	"time"
)

// InsertHost inserts a host into the database
func (m *postgresDBRepo) InsertHost(h models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, query,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		h.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println(err)
		return newID, err
	}

	return newID, nil
}
