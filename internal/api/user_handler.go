package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fsrn12/fitness_tracker_go/internal/store"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userStore store.UserStore
}

func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Created a new User\n")
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		return
	}

	createdUser, err := uh.userStore.CreateUser(&user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create new workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)

}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	paramUserId := chi.URLParam(r, "id")
	if paramUserId == "" {
		http.NotFound(w, r)
		return
	}
	userID, err := strconv.ParseInt(paramUserId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "User ID: %d\n", userID)

	user, err := uh.userStore.GetUserByID(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not find user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func (uh *UserHandler) HandleUpdateUserByID(w http.ResponseWriter, r *http.Request) {

	paramUserId := chi.URLParam(r, "id")
	if paramUserId == "" {
		http.NotFound(w, r)
		return
	}
	userID, err := strconv.ParseInt(paramUserId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "User ID: %d\n", userID)

	existingUser, err := uh.userStore.GetUserByID(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not find user", http.StatusNotFound)
		return
	}

	if existingUser == nil {
		http.NotFound(w, r)
		return
	}

	var updateUserRequest struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Bio      *string `json:"bio"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateUserRequest.Username != nil {
		existingUser.Username = *updateUserRequest.Username
	}
	if updateUserRequest.Email != nil {
		existingUser.Email = *updateUserRequest.Email
	}
	if updateUserRequest.Bio != nil {
		existingUser.Bio = *updateUserRequest.Bio
	}

	err = uh.userStore.UpdateUser(existingUser)
	if err != nil {
		fmt.Println("Failed to update user", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingUser)

}

func (uh *UserHandler) HandleDeleteUserByID(w http.ResponseWriter, r *http.Request) {
	paramUserId := chi.URLParam(r, "id")
	if paramUserId == "" {
		http.NotFound(w, r)
		return
	}

	userID, err := strconv.ParseInt(paramUserId, 10, 64)
	if err != nil {
		http.NotFound(w, r)
	}

	err = uh.userStore.DeleteUser(userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Could not delete user", http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println("User Delete Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
