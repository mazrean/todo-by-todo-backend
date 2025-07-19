package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository"
)

type User struct {
	userRepo *repository.UserRepository
}

func NewUser(version Version, repo *repository.UserRepository) *User {
	return &User{
		userRepo: repo,
	}
}

type UserRequest struct {
	Name string
}

func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := u.userRepo.CreateUser(r.Context(), request.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]int64{"user_id": userID})
}
