package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type TodoRepository struct {
	repository *Repository
}

func NewTodoRepository(repository *Repository) *TodoRepository {
	return &TodoRepository{
		repository: repository,
	}
}

func (tr *TodoRepository) CreateTodo(ctx context.Context, userID int64, title string, description *string, completed bool) (int64, error) {
	slog.Info("CreateTodo called", "userID", userID, "title", title, "description", description, "completed", completed)
	var result int64

	id := tr.repository.getNextID()
	now := time.Now()

	todo := Todo{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   completed,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tr.repository.data.Todos = append(tr.repository.data.Todos, todo)
	result = id
	return result, nil
}

func (tr *TodoRepository) GetTodo(ctx context.Context, id int64) (*Todo, error) {
	tr.repository.mutex.RLock()
	defer tr.repository.mutex.RUnlock()

	for _, todo := range tr.repository.data.Todos {
		if todo.ID == id {
			return &todo, nil
		}
	}

	return nil, fmt.Errorf("todo not found")
}

func (tr *TodoRepository) ListTodosByUser(ctx context.Context, userID int64) ([]Todo, error) {
	tr.repository.mutex.RLock()
	defer tr.repository.mutex.RUnlock()

	var userTodos []Todo
	for _, todo := range tr.repository.data.Todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}

	return userTodos, nil
}

func (tr *TodoRepository) ListTodos(ctx context.Context) ([]Todo, error) {
	tr.repository.mutex.RLock()
	defer tr.repository.mutex.RUnlock()

	return tr.repository.data.Todos, nil
}

func (tr *TodoRepository) UpdateTodo(ctx context.Context, id int64, title string, description *string, completed bool) error {
	return tr.repository.Transaction(ctx, func(ctx context.Context) error {
		tr.repository.mutex.Lock()
		defer tr.repository.mutex.Unlock()

		for i, todo := range tr.repository.data.Todos {
			if todo.ID == id {
				tr.repository.data.Todos[i].Title = title
				tr.repository.data.Todos[i].Description = description
				tr.repository.data.Todos[i].Completed = completed
				tr.repository.data.Todos[i].UpdatedAt = time.Now()
				return nil
			}
		}

		return fmt.Errorf("todo not found")
	})
}

func (tr *TodoRepository) DeleteTodo(ctx context.Context, id int64) error {
	return tr.repository.Transaction(ctx, func(ctx context.Context) error {
		tr.repository.mutex.Lock()
		defer tr.repository.mutex.Unlock()

		for i, todo := range tr.repository.data.Todos {
			if todo.ID == id {
				tr.repository.data.Todos = append(tr.repository.data.Todos[:i], tr.repository.data.Todos[i+1:]...)
				return nil
			}
		}

		return fmt.Errorf("todo not found")
	})
}

func (tr *TodoRepository) MarkTodoCompleted(ctx context.Context, id int64) error {
	return tr.repository.Transaction(ctx, func(ctx context.Context) error {
		tr.repository.mutex.Lock()
		defer tr.repository.mutex.Unlock()

		for i, todo := range tr.repository.data.Todos {
			if todo.ID == id {
				tr.repository.data.Todos[i].Completed = true
				tr.repository.data.Todos[i].UpdatedAt = time.Now()
				return nil
			}
		}

		return fmt.Errorf("todo not found")
	})
}

func (tr *TodoRepository) MarkTodoIncomplete(ctx context.Context, id int64) error {
	return tr.repository.Transaction(ctx, func(ctx context.Context) error {
		tr.repository.mutex.Lock()
		defer tr.repository.mutex.Unlock()

		for i, todo := range tr.repository.data.Todos {
			if todo.ID == id {
				tr.repository.data.Todos[i].Completed = false
				tr.repository.data.Todos[i].UpdatedAt = time.Now()
				return nil
			}
		}

		return fmt.Errorf("todo not found")
	})
}
