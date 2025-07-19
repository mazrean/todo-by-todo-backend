package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

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

func (t *Todo) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func (t *Todo) GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetTodoList, %q", html.EscapeString(r.URL.Path))

	t.todoRepo.ListTodos(r.Context())
}

func (t *Todo) PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var request TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// userID int64, title string, description *string, completed bool
	t.todoRepo.CreateTodo(r.Context(), request.UserID, request.Title, request.Description, request.Completed)
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
	t.todoRepo.UpdateTodo(r.Context(), idInt64, request.Title, request.Description, request.Completed)
}

func (t *Todo) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt64, err := utils.ParseInt64(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	t.todoRepo.DeleteTodo(r.Context(), idInt64)
}
