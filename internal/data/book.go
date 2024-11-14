package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tchenbz/AWT_Test3/internal/validator"
)

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author   	  string    `json:"author"`
	Genre      	  string    `json:"genre"`
	AverageRating float32   `json:"average_rating"`
	CreatedAt     time.Time `json:"-"`
	Version       int32     `json:"version"`
}

type BookModel struct {
	DB *sql.DB
}

func (m BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books (title, author, genre)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []interface{}{book.Title, book.Author, book.Genre}
	return m.DB.QueryRowContext(context.Background(), query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(book.Author != "", "author", "must be provided")}

func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, author, genre, average_rating, created_at, version
		FROM books
		WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Genre,
		&book.AverageRating,
		&book.CreatedAt,
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, genre = $3, average_rating = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []interface{}{
		book.Title,
		book.Author,
		book.Genre,
		book.AverageRating,
		book.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM books
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m BookModel) GetAll(title, author string, filters Filters) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, title, author, genre, average_rating, created_at, version
		FROM books
		WHERE (title ILIKE $1 OR $1 = '')
		AND (author ILIKE $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		"%" + title + "%",
		"%" + author + "%",
		filters.limit(),
		filters.offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	books := []*Book{}

	for rows.Next() {
		var book Book
		err := rows.Scan(
			&totalRecords,
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.AverageRating,
			&book.CreatedAt,
			&book.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return books, metadata, nil
}
