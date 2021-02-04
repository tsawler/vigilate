package dbrepo

import (
	"context"
	"time"
)

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
