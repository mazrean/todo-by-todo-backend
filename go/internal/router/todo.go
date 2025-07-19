package router

import (
	"fmt"
	"html"
	"net/http"
)

type Todo struct{}

type Version string

func NewTodo(version Version) *Todo {
	return &Todo{}
}

func (e *Todo) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func (e *Todo) GetTodoList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetTodoList, %q", html.EscapeString(r.URL.Path))
}

func (e *Todo) PostTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PostTodo, %q", html.EscapeString(r.URL.Path))
}

func (e *Todo) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UpdateTodo, %q", html.EscapeString(r.URL.Path))
}

func (e *Todo) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DeleteTodo, %q", html.EscapeString(r.URL.Path))
}
