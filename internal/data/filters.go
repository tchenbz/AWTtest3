package data

import (
	"strings"

	"github.com/tchenbz/AWT_Test1/internal/validator"
)

// The Filters type will contain the fields related to pagination
// and eventually the fields related to sorting.
type Filters struct {
	Page         int      // Page number requested by the client
	PageSize     int      // Number of records per page
	Sort         string   // Sort column and direction (e.g., "id" or "-id")
	SortSafeList []string // Allowed fields for sorting to prevent SQL injection
}

// Metadata to be included in responses
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// ValidateFilters validates pagination and sorting fields in Filters.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 500, "page", "must not exceed 500")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// Set a limit on the number of records per page
func (f Filters) limit() int {
	return f.PageSize
}

// Set the offset to retrieve records for a specific page
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// Sort column for SQL queries, ensuring safety by restricting to allowed fields.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	// If an invalid sort value is provided, prevent execution for security
	panic("unsafe sort parameter: " + f.Sort)
}

// Sort direction for SQL queries, either ASC or DESC based on "-" prefix.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// Calculate metadata for pagination
func calculateMetaData(totalRecords int, currentPage int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
