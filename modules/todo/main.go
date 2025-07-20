package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	todoapi "github.com/mazrean/todo-by-todo-backend/modules/todo/internal/bindings/todo/api/todo-api"
	"github.com/mazrean/todo-by-todo-backend/modules/todo/internal/di"
	"github.com/mazrean/todo-by-todo-backend/modules/todo/internal/repository"
	"github.com/mazrean/todo-by-todo-backend/modules/todo/internal/repository/config"
)

//go:generate wkg wit build --output bindings.wasm
//go:generate go tool wit-bindgen-go generate --world todo --out internal/bindings bindings.wasm

func init() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: parseLogLevel(logLevel),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "development"
	}

	environment := config.Environment(env)

	dbFilePath, ok := os.LookupEnv("DB_FILEPATH")
	if !ok {
		panic("DB_FILEPATH is not set")
	}

	dbCfg := &repository.JSONConfig{
		DataPath: dbFilePath,
	}

	app, err := di.InjectCLI(
		dbCfg,
		environment,
	)
	if err != nil {
		panic(fmt.Errorf("failed to inject dependencies: %w", err))
	}

	todoapi.Exports.ListTodos = app.Todo.GetTodoListHandler
	todoapi.Exports.CreateTodo = app.Todo.PostTodoHandler
	todoapi.Exports.UpdateTodo = app.Todo.UpdateTodoHandler
	todoapi.Exports.DeleteTodo = app.Todo.DeleteTodoHandler

	todoapi.Exports.CreateUser = app.User.CreateUserHandler
}

func main() {}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
