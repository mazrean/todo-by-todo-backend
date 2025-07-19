package di

//go:generate go tool kessoku $GOFILE

import (
	"github.com/mazrean/kessoku"
	repo "github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository"
	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/router"
)

// NOTE: repoに変数名を上書きする
var _ = kessoku.Inject[*router.Router](
	"InjectCLI",
	kessoku.Provide(router.NewRouter),
	kessoku.Provide(router.NewTodo),
	kessoku.Provide(repo.NewTodoRepository),
	kessoku.Provide(repo.NewDB),
)
