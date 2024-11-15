package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tchenbz/AWTtest3/internal/data"
	"github.com/tchenbz/AWTtest3/internal/validator"
)

func (a *applicationDependencies) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title        string `json:"title"`
		Author string `json:"author"`
		Genre    string `json:"genre"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:    input.Title,
		Author:   input.Author,
		Genre:    input.Genre,
	}

	v := validator.New()
	data.ValidateBook(v, book)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.bookModel.Insert(book)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	book, err := a.bookModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	book, err := a.bookModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title        *string `json:"title"`
		Author 		 *string `json:"author"`
		Genre    	 *string `json:"genre"`
	}

	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Genre != nil {
		book.Genre = *input.Genre
	}

	v := validator.New()
	data.ValidateBook(v, book)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.bookModel.Update(book)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{"book": book}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.bookModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"message": "book successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string
		Author string
		data.Filters
	}

	query := r.URL.Query()
	input.Title = a.getSingleQueryParameter(query, "title", "")
	input.Author = a.getSingleQueryParameter(query, "author", "")
	input.Filters.Page = a.getSingleIntegerParameter(query, "page", 1, validator.New())
	input.Filters.PageSize = a.getSingleIntegerParameter(query, "page_size", 10, validator.New())
	input.Filters.Sort = a.getSingleQueryParameter(query, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "Author", "-id", "-title", "-author"}

	v := validator.New()
	data.ValidateFilters(v, input.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := a.bookModel.GetAll(input.Title, input.Author, input.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"books": books,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
