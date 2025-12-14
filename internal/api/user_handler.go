package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/fsrn12/fitness_tracker_go/internal/store"
	"github.com/fsrn12/fitness_tracker_go/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateUserRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" || len(req.Username) < 2 || len(req.Username) > 50 {
		return errors.New("user name is required. It must be at least 2 to 50 characters long")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required and it must be between 8 to 32 characters long")
	}

	_, err := utils.ValidatePassword(req.Password)
	if err != nil {
		return err
	}

	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding new user register request %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = uh.validateUserRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	_, err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: registering user %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": user})

}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.logger.Printf("ERROR - CreateUserRequestDecoding: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdUser, err := uh.userStore.CreateUser(&user)
	if err != nil {
		uh.logger.Printf("ERROR - CreateUser(): %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create user"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": createdUser})
}

func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: GetParamID => %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	user, err := uh.userStore.GetUserByID(userID)
	if err != nil {
		uh.logger.Printf("ERROR - GetWorkoutByID(): %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": user})
}

func (uh *UserHandler) HandleUpdateUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: GetParamID => %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user id"})
		return
	}

	existingUser, err := uh.userStore.GetUserByID(userID)
	if err != nil {
		uh.logger.Printf("ERROR - GetUserByID() -> UpdateUser(): %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		uh.logger.Printf("ERROR - updateUserRequestDecoding: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
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
		uh.logger.Printf("ERROR - UpdateUser(): %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update new user"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": existingUser})
}

func (uh *UserHandler) HandleDeleteUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: GetParamID => %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user id"})
		return
	}

	err = uh.userStore.DeleteUser(userID)
	if err == sql.ErrNoRows {
		uh.logger.Printf("ERROR - DeleteUser(): %v\n", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "no record found"})
		return
	}

	if err != nil {
		uh.logger.Printf("ERROR - DeleteUser(): %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// w.WriteHeader(http.StatusNoContent)
	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"data": "user deleted successfully"})
}

func (uh *UserHandler) HandleGetUserToken(w http.ResponseWriter, r *http.Request) {

}
