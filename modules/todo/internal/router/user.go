package router

import (
	"context"
	"log/slog"

	todoapi "github.com/mazrean/todo-by-todo-backend/modules/todo/internal/bindings/todo/api/todo-api"
	"github.com/mazrean/todo-by-todo-backend/modules/todo/internal/repository"
	"go.bytecodealliance.org/cm"
)

type User struct {
	userRepo *repository.UserRepository
}

func NewUser(repo *repository.UserRepository) *User {
	return &User{
		userRepo: repo,
	}
}

type UserRequest struct {
	Name string
}

func (u *User) CreateUserHandler(request todoapi.UserRequest) (result cm.Option[todoapi.APIError]) {
	ctx := context.Background()

	_, err := u.userRepo.CreateUser(ctx, request.Name)
	if err != nil {
		slog.Error("failed to delete todo", "error", err)
		return cm.Some(todoapi.APIErrorInternalError("Failed to delete todo"))
	}

	return cm.None[todoapi.APIError]()
}
