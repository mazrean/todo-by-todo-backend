package router

import (
	"net/http"
)

type Router struct {
	addr    string
	example *Example
}

func NewRouter(
	addr Addr,
	example *Example,
) *Router {
	return &Router{
		addr:    string(addr),
		example: example,
	}
}

type Addr string

func (r *Router) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/example", http.HandlerFunc(r.example.Handler))

	return http.ListenAndServe(
		r.addr,
		mux,
	)
}
