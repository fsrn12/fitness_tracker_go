package api

import (
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
