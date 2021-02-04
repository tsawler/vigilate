package dbrepo

import (
	"context"
	"database/sql"
	"github.com/tsawler/vigilate/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// AllUsers returns all users
func (m *postgresDBRepo) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, last_name, first_name, email, user_active, created_at, updated_at FROM users`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(&s.ID, &s.LastName, &s.FirstName, &s.Email, &s.UserActive, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		// Append it to the slice
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// GetUserById returns a user by id
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, first_name, last_name,  user_active, access_level, email, 
			created_at, updated_at, use_tfa, avatar
			FROM users where id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.UserActive,
		&u.AccessLevel,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.UseTfa,
		&u.Avatar,
	)

	if err != nil {
		log.Println(err)
		return u, err
	}

	return u, nil
}

// Authenticate authenticates
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	var userActive int
	row := m.DB.QueryRowContext(ctx, "SELECT id, password, user_active FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword, &userActive)
	if err == sql.ErrNoRows {
		return 0, "", models.ErrInvalidCredentials
	} else if err != nil {
		log.Println(err)
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", models.ErrInvalidCredentials
	} else if err != nil {
		log.Println(err)
		return 0, "", err
	}

	if userActive == 0 {
		return 0, "", models.ErrInactiveAccount
	}

	// Otherwise, the password is correct. Return the user ID and hashed password.
	return id, hashedPassword, nil
}

// Insert inserts a token into remember_tokens for a user
func (m *postgresDBRepo) InsertRememberMeToken(id int, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "insert into remember_tokens (user_id, remember_token) values ($1, $2)"
	_, err := m.DB.ExecContext(ctx, stmt, id, token)
	if err != nil {
		return err
	}
	return nil
}

// DeleteToken deletes a remember me token
func (m *postgresDBRepo) DeleteToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "delete from remember_tokens where remember_token = $1"
	_, err := m.DB.ExecContext(ctx, stmt, token)
	if err != nil {
		return err
	}

	return nil
}

// CheckForToken checks for a valid remember me token
func (m *postgresDBRepo) CheckForToken(id int, token string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "SELECT id  FROM remember_tokens where user_id = $1 and remember_token = $2"
	row := m.DB.QueryRowContext(ctx, stmt, id, token)
	err := row.Scan(&id)
	return err == nil
}
