package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tchenbz/AWT_Test1/internal/data"
	"github.com/tchenbz/AWT_Test1/internal/validator"
)

func (a *applicationDependencies) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	product := &data.Product{
		Name:        input.Name,
		Description: input.Description,
		Category:    input.Category,
		ImageURL:    input.ImageURL,
	}

	v := validator.New()
	data.ValidateProduct(v, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.productModel.Insert(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/%d", product.ID))

	data := envelope{"product": product}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayProductHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Fetch the product from the database
	product, err := a.productModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send the product data in JSON format
	data := envelope{"product": product}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Retrieve the existing product from the database
	product, err := a.productModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Create a temporary struct for incoming updates
	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		ImageURL    *string `json:"image_url"`
	}

	// Decode the request JSON into the input struct
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Update the product fields as needed
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.ImageURL != nil {
		product.ImageURL = *input.ImageURL
	}

	// Validate the updated product
	v := validator.New()
	data.ValidateProduct(v, product)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Perform the update in the database
	err = a.productModel.Update(product)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond with the updated product data
	data := envelope{"product": product}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Attempt to delete the product
	err = a.productModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a JSON response confirming deletion
	data := envelope{"message": "product successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a struct to hold the query parameters
	var input struct {
		Name     string
		Category string
		data.Filters
	}

	// Load the query parameters into the input struct
	query := r.URL.Query()
	input.Name = a.getSingleQueryParameter(query, "name", "")
	input.Category = a.getSingleQueryParameter(query, "category", "")
	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "name", "category", "-id", "-name", "-category"}

	// Validate the filters
	v := validator.New()
	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the products based on the filters
	products, metadata, err := a.productModel.GetAll(input.Name, input.Category, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send the product list and metadata in JSON format
	data := envelope{
		"products": products,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
