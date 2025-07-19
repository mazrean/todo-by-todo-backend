package router

import (
	"fmt"
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
	basePath := "/api"

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc(fmt.Sprintf("GET %s/todos", basePath), r.todo.GetTodoListHandler)
	mux.HandleFunc(fmt.Sprintf("POST %s/todos", basePath), r.todo.PostTodoHandler)
	mux.HandleFunc(fmt.Sprintf("PUT %s/todos/{id}", basePath), r.todo.UpdateTodoHandler)
	mux.HandleFunc(fmt.Sprintf("DELETE %s/todos/{id}", basePath), r.todo.DeleteTodoHandler)

	mux.HandleFunc(fmt.Sprintf("POST %s/users", basePath), r.user.CreateUserHandler)

	return http.ListenAndServe(
		r.addr,
		mux,
	)
}
