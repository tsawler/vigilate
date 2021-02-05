package dbrepo

import (
	"context"
	"github.com/tsawler/vigilate/internal/models"
	"log"
	"time"
)

// AllPreferences returns a slice of preferences
func (m *postgresDBRepo) AllPreferences() ([]models.Preference, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "SELECT id, name, preference FROM preferences"

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []models.Preference

	for rows.Next() {
		s := &models.Preference{}
		err = rows.Scan(&s.ID, &s.Name, &s.Preference)
		if err != nil {
			return nil, err
		}
		preferences = append(preferences, *s)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return preferences, nil
}

// SetSystemPref updates a system preference setting
func (m *postgresDBRepo) SetSystemPref(name, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from preferences where name = $1`
	_, _ = m.DB.ExecContext(ctx, stmt, name)

	query := `
		INSERT INTO preferences (
			  	name, preference, created_at, updated_at
			  ) VALUES ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, query, name, value, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// InsertOrUpdateSitePreferences inserts or updates all site prefs from map
func (m *postgresDBRepo) InsertOrUpdateSitePreferences(pm map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for k, v := range pm {
		query := `delete from preferences where name = $1`

		_, err := m.DB.ExecContext(ctx, query, k)
		if err != nil {
			return err
		}

		query = `insert into preferences (name, preference, created_at, updated_at)
			values ($1, $2, $3, $4)`

		_, err = m.DB.ExecContext(ctx, query, k, v, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
