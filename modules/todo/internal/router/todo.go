package router

import (
	"context"
	"fmt"
	"log/slog"

	todoapi "github.com/mazrean/todo-by-todo-backend/modules/todo/internal/bindings/todo/api/todo-api"
	"github.com/mazrean/todo-by-todo-backend/modules/todo/internal/repository"
	"go.bytecodealliance.org/cm"
)

type Todo struct {
	todoRepo *repository.TodoRepository
}

func NewTodo(repo *repository.TodoRepository) *Todo {
	return &Todo{
		todoRepo: repo,
	}
}

func (t *Todo) GetTodoListHandler() (result todoapi.TodosResult) {
	slog.Info("GetTodoListHandler called")
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in CreateUserHandler", "panic", r)
			result.Error = cm.Some(todoapi.APIErrorInternalError(fmt.Sprintf("panic: %v", r)))
			return
		}
	}()
	ctx := context.Background()

	todos, err := t.todoRepo.ListTodos(ctx)
	if err != nil {
		slog.Error("failed to list todos", "error", err)
		return todoapi.TodosResult{
			Error: cm.Some(todoapi.APIErrorInternalError("Failed to list todos")),
		}
	}

	resTodos := make([]todoapi.Todo, 0, len(todos))
	for _, todo := range todos {
		var description cm.Option[string]
		if todo.Description != nil {
			description = cm.Some(*todo.Description)
		} else {
			description = cm.None[string]()
		}

		var completed cm.Option[bool]
		completed = cm.Some(todo.Completed)

		resTodo := todoapi.Todo{
			ID:          todoapi.TodoID(todo.ID),
			UserID:      todoapi.UserID(todo.UserID),
			Title:       todo.Title,
			Description: description,
			Completed:   completed,
		}
		resTodos = append(resTodos, resTodo)
	}

	return todoapi.TodosResult{
		Todos: cm.ToList(resTodos),
		Error: cm.None[todoapi.APIError](),
	}
}

func (t *Todo) PostTodoHandler(request todoapi.TodoRequest) (result todoapi.CreateResult) {
	slog.Info("PostTodoHandler called", "completed", request.Completed, "title", request.Title, "description", request.Description, "userID", request.UserID, "request", request)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in PostTodoHandler", "panic", r)
			result.Error = cm.Some(todoapi.APIErrorInternalError(fmt.Sprintf("panic: %v", r)))
			return
		}
	}()
	ctx := context.Background()

	_, err := t.todoRepo.CreateTodo(
		ctx,
		int64(request.UserID),
		request.Title,
		request.Description.Some(),
		request.Completed,
	)
	if err != nil {
		slog.Error("failed to create todo", "error", err)
		return todoapi.CreateResult{
			Error: cm.Some(todoapi.APIErrorInternalError("Failed to create todo")),
		}
	}

	return todoapi.CreateResult{
		Error: cm.None[todoapi.APIError](),
	}
}

func (t *Todo) UpdateTodoHandler(id todoapi.TodoID, request todoapi.TodoRequest) (result cm.Option[todoapi.APIError]) {
	slog.Info("UpdateTodoHandler called", "id", id, "request", request)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in UpdateTodoHandler", "panic", r)
			result = cm.Some(todoapi.APIErrorInternalError(fmt.Sprintf("panic: %v", r)))
			return
		}
	}()
	ctx := context.Background()

	err := t.todoRepo.UpdateTodo(ctx, int64(id), request.Title, request.Description.Some(), request.Completed)
	if err != nil {
		slog.Error("failed to update todo", "error", err)
		return cm.Some(todoapi.APIErrorInternalError("Failed to update todo"))
	}

	return cm.None[todoapi.APIError]()
}

func (t *Todo) DeleteTodoHandler(id todoapi.TodoID) (result cm.Option[todoapi.APIError]) {
	slog.Info("DeleteTodoHandler called", "id", id)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in DeleteTodoHandler", "panic", r)
			result = cm.Some(todoapi.APIErrorInternalError(fmt.Sprintf("panic: %v", r)))
			return
		}
	}()
	ctx := context.Background()

	err := t.todoRepo.DeleteTodo(ctx, int64(id))
	if err != nil {
		slog.Error("failed to delete todo", "error", err)
		return cm.Some(todoapi.APIErrorInternalError("Failed to delete todo"))
	}

	return cm.None[todoapi.APIError]()
}
