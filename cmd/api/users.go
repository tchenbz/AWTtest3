// cmd/api/user.go
package main

import (
	"net/http"
	"fmt"
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
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))  // Set the Location header to the new user's URL

	err = a.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// showUserHandler handles retrieving a user by ID.
func (a *applicationDependencies) showUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	user, err := a.userModel.Get(id)
	if err != nil {
		if err.Error() == "user not found" {
			a.notFoundResponse(w, r)
		} else {
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	err = a.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// updateUserHandler handles updating a user's information.
func (a *applicationDependencies) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	user, err := a.userModel.Get(id)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	var input struct {
		Email    *string `json:"email"`
		FullName *string `json:"full_name"`
	}
	err = a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = a.userModel.Update(user)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	err = a.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// deleteUserHandler handles deleting a user by ID.
func (a *applicationDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}
	err = a.userModel.Delete(id)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	err = a.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

