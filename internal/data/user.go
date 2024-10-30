// internal/data/user.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"github.com/tchenbz/AWT_Quiz3/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

type UserModel struct {
	DB *sql.DB
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Email != "", "email", "must be provided")
	v.Check(user.FullName != "", "full_name", "must be provided")
	v.Check(len(user.FullName) <= 100, "full_name", "must not be more than 100 characters long")
}

// Insert adds a new user to the database
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (email, full_name)
		VALUES ($1, $2)
		RETURNING id, created_at, version`
	args := []any{user.Email, user.FullName}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

// Get retrieves a specific user by ID from the database
func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, errors.New("invalid user ID")
	}
	query := `
		SELECT id, email, full_name, created_at, version
		FROM users
		WHERE id = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.FullName, &user.CreatedAt, &user.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update modifies an existing user's information in the database
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, version = version + 1
		WHERE id = $3
		RETURNING version`
	args := []any{user.Email, user.FullName, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
}

// Delete removes a user by ID from the database
func (m UserModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("invalid user ID")
	}
	query := `
		DELETE FROM users
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

