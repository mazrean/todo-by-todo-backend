package router

import (
	"net/http"
)

type Router struct {
	addr string
	todo *Todo
}

func NewRouter(
	addr Addr,
	todo *Todo,
) *Router {
	return &Router{
		addr: string(addr),
		todo: todo,
	}
}

type Addr string

func (r *Router) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /todos", r.todo.GetTodoListHandler)
	mux.HandleFunc("POST /todos", r.todo.PostTodoHandler)
	mux.HandleFunc("PUT /todos/{id}", r.todo.UpdateTodoHandler)
	mux.HandleFunc("DELETE /todos/{id}", r.todo.DeleteTodoHandler)

	return http.ListenAndServe(
		r.addr,
		mux,
	)
}
