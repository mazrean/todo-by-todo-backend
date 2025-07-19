package di

//go:generate go tool kessoku $GOFILE

import (
	"github.com/mazrean/kessoku"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/router"
)

var _ = kessoku.Inject[*router.Router](
	"InjectCLI",
	kessoku.Provide(router.NewRouter),
	kessoku.Provide(router.NewExample),
	kessoku.Provide(repository.NewDB),
)
