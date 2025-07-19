package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/mazrean/todo-by-todo-backend/modules/todo/internal/repository/config"
)

type Repository struct {
	dataPath string
	data     *JSONData
	mutex    sync.RWMutex
}

type JSONConfig struct {
	DataPath string
}

func NewJSON(jsonConf *JSONConfig, env config.Environment) (*Repository, error) {
	absPath, err := filepath.Abs(jsonConf.DataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	repo := &Repository{
		dataPath: absPath,
		data:     &JSONData{NextID: 1},
	}

	if err := repo.loadData(); err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	return repo, nil
}

func (r *Repository) loadData() error {
	data, err := os.ReadFile(r.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return r.saveData()
		}
		return fmt.Errorf("failed to read data file: %w", err)
	}

	if err := json.Unmarshal(data, r.data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

func (r *Repository) saveData() error {
	dir := filepath.Dir(r.dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(r.dataPath, data, fs.FileMode(0644)); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

func (r *Repository) Close() error {
	return nil
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := fn(ctx); err != nil {
		return err
	}

	return r.saveData()
}

func (r *Repository) getNextID() int64 {
	id := r.data.NextID
	r.data.NextID++
	return id
}
