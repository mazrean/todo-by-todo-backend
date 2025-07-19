package repository

import (
	"context"
	"database/sql"

	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository/db"
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
	queries := tr.repository.GetQueriesWithTx(ctx)

	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}

	result, err := queries.CreateTodo(ctx, db.CreateTodoParams{
		UserID:      userID,
		Title:       title,
		Description: desc,
		Completed:   sql.NullBool{Bool: completed, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tr *TodoRepository) GetTodo(ctx context.Context, id int64) (*db.Todo, error) {
	queries := tr.repository.GetQueriesWithTx(ctx)

	todo, err := queries.GetTodo(ctx, id)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (tr *TodoRepository) ListTodosByUser(ctx context.Context, userID int64) ([]db.Todo, error) {
	queries := tr.repository.GetQueriesWithTx(ctx)

	return queries.ListTodosByUser(ctx, userID)
}

func (tr *TodoRepository) ListTodos(ctx context.Context) ([]db.Todo, error) {
	queries := tr.repository.GetQueriesWithTx(ctx)

	return queries.ListTodos(ctx)
}

func (tr *TodoRepository) UpdateTodo(ctx context.Context, id int64, title string, description *string, completed bool) error {
	queries := tr.repository.GetQueriesWithTx(ctx)

	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}

	return queries.UpdateTodo(ctx, db.UpdateTodoParams{
		Title:       title,
		Description: desc,
		Completed:   sql.NullBool{Bool: completed, Valid: true},
		ID:          id,
	})
}

func (tr *TodoRepository) DeleteTodo(ctx context.Context, id int64) error {
	queries := tr.repository.GetQueriesWithTx(ctx)

	return queries.DeleteTodo(ctx, id)
}

func (tr *TodoRepository) MarkTodoCompleted(ctx context.Context, id int64) error {
	queries := tr.repository.GetQueriesWithTx(ctx)

	return queries.MarkTodoCompleted(ctx, id)
}

func (tr *TodoRepository) MarkTodoIncomplete(ctx context.Context, id int64) error {
	queries := tr.repository.GetQueriesWithTx(ctx)

	return queries.MarkTodoIncomplete(ctx, id)
}
