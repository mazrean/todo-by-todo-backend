package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository"
)

type Todo struct{
	todoRepo *repository.TodoRepository
}

type Version string

func NewTodo(version Version, repo *repository.TodoRepository) *Todo {
	return &Todo{
		todoRepo: repo,
	}
}

type TodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (e *Todo) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func (e *Todo) GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetTodoList, %q", html.EscapeString(r.URL.Path))
}

func (e *Todo) PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var request TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "PostTodo: title=%s, done=%v", request.Title, request.Done)
}

func (e *Todo) UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var request TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "UpdateTodo id=%s: title=%s, done=%v", id, request.Title, request.Done)
}

func (e *Todo) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "DeleteTodo %s, %q", id, html.EscapeString(r.URL.Path))
}
