package router

import (
	"net/http"
)

type Router struct {
	addr string
	todo *Todo
	user *User
}

func NewRouter(
	addr Addr,
	todo *Todo,
	user *User,
) *Router {
	return &Router{
		addr: string(addr),
		todo: todo,
		user: user,
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

	mux.HandleFunc("POST /users", r.user.CreateUserHandler)

	return http.ListenAndServe(
		r.addr,
		mux,
	)
}
