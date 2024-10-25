package main

import (
	"fmt"
	"net/http"
	"github.com/tchenbz/AWT_Quiz3/internal/data"
	"github.com/tchenbz/AWT_Quiz3/internal/validator"
)

func (a *applicationDependencies) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		FullName string `json:"full_name"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Email:    input.Email,
		FullName: input.FullName,
	}

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.userModel.Insert(user)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

	err = a.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
