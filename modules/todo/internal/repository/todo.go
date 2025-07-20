package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
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

	bytes, err := json.MarshalIndent(tr.repository.data, "", "  ")
	if err != nil {
		return 0, err
	}

	// ファイルに書き込み
	if err := os.WriteFile(tr.repository.dataPath, bytes, 0o644); err != nil {
		return 0, err
	}
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

	// JSON ファイルを読み込む
	bytes, err := os.ReadFile(tr.repository.dataPath)
	if err != nil {
		return nil, err
	}

	// 一時的な変数へ Unmarshal
	var j JSONData
	if err := json.Unmarshal(bytes, &j); err != nil {
		return nil, err
	}

	return j.Todos, nil
}

func (tr *TodoRepository) UpdateTodo(
	ctx context.Context,
	id int64,
	title string,
	description *string,
	completed bool,
) error {
	// メモリ／ファイル双方の更新は排他制御
	tr.repository.mutex.Lock()
	defer tr.repository.mutex.Unlock()

	// 対象 Todo を検索して更新
	for i, todo := range tr.repository.data.Todos {
		if todo.ID == id {
			tr.repository.data.Todos[i].Title = title
			tr.repository.data.Todos[i].Description = description
			tr.repository.data.Todos[i].Completed = completed
			tr.repository.data.Todos[i].UpdatedAt = time.Now()

			// メモリ上のデータを JSON にシリアライズ
			bytes, err := json.MarshalIndent(tr.repository.data, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			// JSON ファイルに上書き
			if err := os.WriteFile(tr.repository.dataPath, bytes, 0o644); err != nil {
				return fmt.Errorf("failed to write JSON file: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("todo not found: id=%d", id)
}

func (tr *TodoRepository) DeleteTodo(
	ctx context.Context,
	id int64,
) error {

	// Todo を検索し、見つかればスライスから削除
	for i, todo := range tr.repository.data.Todos {
		if todo.ID == id {
			tr.repository.data.Todos = append(
				tr.repository.data.Todos[:i],
				tr.repository.data.Todos[i+1:]...,
			)

			// 削除後の全データを JSON にシリアライズ
			bytes, err := json.MarshalIndent(tr.repository.data, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}

			// JSON ファイルに上書き
			if err := os.WriteFile(tr.repository.dataPath, bytes, 0o644); err != nil {
				return fmt.Errorf("failed to write JSON file: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("todo not found: id=%d", id)
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
