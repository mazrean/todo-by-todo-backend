package router

import (
	"net/http"
)

type Router struct {
	addr    string
	todo    *Todo
}

func NewRouter(
	addr Addr,
	todo *Todo,
) *Router {
	return &Router{
		addr:    string(addr),
		todo: todo,
	}
}

type Addr string

func (r *Router) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", r.todo.GetTodoListHandler)
	mux.HandleFunc("POST /todo", r.todo.PostTodoHandler)
	mux.HandleFunc("PUT /todo/{id}", r.todo.UpdateTodoHandler)
	mux.HandleFunc("DELETE /todo/{id}", r.todo.DeleteTodoHandler)

	return http.ListenAndServe(
		r.addr,
		mux,
	)
}
