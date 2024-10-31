package validator

import (
	"slices"
)
 
type Validator struct {
    Errors map[string]string
} 

func New() *Validator {
    return &Validator {
        Errors: make(map[string]string),
    }
}

func (v *Validator) IsEmpty() bool {
    return len(v.Errors) == 0
}

func (v *Validator) AddError(key string, message string) {
    _, exists := v.Errors[key]
    if !exists {
        v.Errors[key] = message
    }
}

func (v *Validator) Check(acceptable bool, key string, message string) {
    if !acceptable {
       v.AddError(key, message)
    }
}

func PermittedValue(value string, permittedValues ...string) bool {
	return slices.Contains(permittedValues, value)
}

// ValidateProduct validates the fields of a product.
// func ValidateProduct(v *Validator, product *Product) {
// 	v.Check(product.Name != "", "name", "must be provided")
// 	v.Check(len(product.Name) <= 100, "name", "must not exceed 100 characters")
// 	v.Check(product.Category != "", "category", "must be provided")
// 	v.Check(len(product.Category) <= 50, "category", "must not exceed 50 characters")
// 	v.Check(product.ImageURL != "", "image_url", "must be provided")
// 	v.Check(len(product.ImageURL) <= 255, "image_url", "must not exceed 255 characters")
// }

// // ValidateReview validates the fields of a review.
// func ValidateReview(v *Validator, review *Review) {
// 	v.Check(review.Content != "", "content", "must be provided")
// 	v.Check(len(review.Content) <= 500, "content", "must not exceed 500 characters")
// 	v.Check(review.Author != "", "author", "must be provided")
// 	v.Check(len(review.Author) <= 50, "author", "must not exceed 50 characters")
// 	v.Check(review.Rating >= 1 && review.Rating <= 5, "rating", "must be between 1 and 5")
// 	v.Check(review.HelpfulCount >= 0, "helpful_count", "must be a positive integer")
// }