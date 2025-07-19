package repository

import (
	"context"

	"github.com/mazrean/mazrean/todo-by-todo-backend/internal/repository/db"
)

type UserRepository struct {
	repository *Repository
}

func NewUserRepository(repository *Repository) *UserRepository {
	return &UserRepository{
		repository: repository,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, name string) (int64, error) {
	queries := ur.repository.GetQueriesWithTx(ctx)
	
	result, err := queries.CreateUser(ctx, name)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

func (ur *UserRepository) GetUser(ctx context.Context, id int64) (*db.User, error) {
	queries := ur.repository.GetQueriesWithTx(ctx)
	
	user, err := queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, id int64, name string) error {
	queries := ur.repository.GetQueriesWithTx(ctx)
	
	return queries.UpdateUser(ctx, db.UpdateUserParams{
		Name: name,
		ID:   id,
	})
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	queries := ur.repository.GetQueriesWithTx(ctx)
	
	return queries.DeleteUser(ctx, id)
}

func (ur *UserRepository) ListUsers(ctx context.Context) ([]db.User, error) {
	queries := ur.repository.GetQueriesWithTx(ctx)
	
	return queries.ListUsers(ctx)
}