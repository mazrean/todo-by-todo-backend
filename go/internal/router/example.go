package router

import (
	"fmt"
	"html"
	"net/http"
)

type Example struct{}

type Version string

func NewExample(version Version) *Example {
	return &Example{}
}

func (e *Example) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
