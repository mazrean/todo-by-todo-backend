package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/utils"
)

type Todo struct {
	todoRepo *repository.TodoRepository
}

type Version string

func NewTodo(version Version, repo *repository.TodoRepository) *Todo {
	return &Todo{
		todoRepo: repo,
	}
}

type TodoRequest struct {
	UserID      int64   `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Completed   bool    `json:"completed"`
}

type TodoResponse struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Todo) GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := t.todoRepo.ListTodos(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list todos: %v", err), http.StatusInternalServerError)
		return
	}

	var response []TodoResponse
	for _, todo := range todos {
		response = append(response, TodoResponse{
			ID:          todo.ID,
			UserID:      todo.UserID,
			Title:       todo.Title,
			Description: &todo.Description.String,
			Completed:   todo.Completed.Bool,
			CreatedAt:   todo.CreatedAt.Time,
			UpdatedAt:   todo.UpdatedAt.Time,
		})
	}
	WriteJSON(w, http.StatusOK, response)
}

func (t *Todo) PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var request TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := t.todoRepo.CreateTodo(r.Context(), request.UserID, request.Title, request.Description, request.Completed)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create todo: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (t *Todo) UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt64, err := utils.ParseInt64(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var request TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = t.todoRepo.UpdateTodo(r.Context(), idInt64, request.Title, request.Description, request.Completed)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (t *Todo) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt64, err := utils.ParseInt64(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = t.todoRepo.DeleteTodo(r.Context(), idInt64)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
