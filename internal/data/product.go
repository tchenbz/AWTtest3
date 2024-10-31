package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/tchenbz/AWT_Test1/internal/validator"
)

type Product struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	ImageURL      string    `json:"image_url"`
	AverageRating float32   `json:"average_rating"`
	CreatedAt     time.Time `json:"-"`
	Version       int32     `json:"version"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) Insert(product *Product) error {
	query := `
		INSERT INTO products (name, description, category, image_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{product.Name, product.Description, product.Category, product.ImageURL}
	return m.DB.QueryRowContext(context.Background(), query, args...).Scan(&product.ID, &product.CreatedAt, &product.Version)
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(product.Category != "", "category", "must be provided")
	v.Check(product.ImageURL != "", "image_url", "must be provided")
}

// Get retrieves a specific product by ID.
func (m ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, description, category, image_url, average_rating, created_at, version
		FROM products
		WHERE id = $1`

	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.ImageURL,
		&product.AverageRating,
		&product.CreatedAt,
		&product.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

// Update modifies a product in the database.
func (m ProductModel) Update(product *Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, category = $3, image_url = $4, average_rating = $5, version = version + 1
		WHERE id = $6
		RETURNING version`

	args := []interface{}{
		product.Name,
		product.Description,
		product.Category,
		product.ImageURL,
		product.AverageRating,
		product.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&product.Version)
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

// Delete removes a product from the database by its ID.
func (m ProductModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM products
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

// GetAll retrieves a list of products with optional filtering, sorting, and pagination.
func (m ProductModel) GetAll(name, category string, filters Filters) ([]*Product, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, description, category, image_url, average_rating, created_at, version
		FROM products
		WHERE (name ILIKE $1 OR $1 = '')
		AND (category ILIKE $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	args := []interface{}{
		"%" + name + "%",
		"%" + category + "%",
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
	products := []*Product{}

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.ImageURL,
			&product.AverageRating,
			&product.CreatedAt,
			&product.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	return products, metadata, nil
}
