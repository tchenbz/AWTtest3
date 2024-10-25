package data

import (
	"context"
	"database/sql"
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

func (c UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (email, full_name)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`
	args := []any{user.Email, user.FullName}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}
