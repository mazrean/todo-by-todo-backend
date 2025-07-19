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

	mux.HandleFunc("/todo", func(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.todo.GetTodoList(w, req)
	case http.MethodPost:
		r.todo.PostTodo(w, req)
	case http.MethodPut:
		r.todo.UpdateTodo(w, req)
	case http.MethodDelete:
		r.todo.DeleteTodo(w, req)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
})

	return http.ListenAndServe(
		r.addr,
		mux,
	)
}
