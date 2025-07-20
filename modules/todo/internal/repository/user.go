package repository

import (
	"context"
	"fmt"
	"time"
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
	var result int64

	ur.repository.mutex.Lock()
	defer ur.repository.mutex.Unlock()

	id := ur.repository.getNextID()
	now := time.Now()

	user := User{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ur.repository.data.Users = append(ur.repository.data.Users, user)
	result = id

	return result, nil
}

func (ur *UserRepository) GetUser(ctx context.Context, id int64) (*User, error) {
	ur.repository.mutex.RLock()
	defer ur.repository.mutex.RUnlock()

	for _, user := range ur.repository.data.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (ur *UserRepository) UpdateUser(ctx context.Context, id int64, name string) error {
	return ur.repository.Transaction(ctx, func(ctx context.Context) error {
		ur.repository.mutex.Lock()
		defer ur.repository.mutex.Unlock()

		for i, user := range ur.repository.data.Users {
			if user.ID == id {
				ur.repository.data.Users[i].Name = name
				ur.repository.data.Users[i].UpdatedAt = time.Now()
				return nil
			}
		}

		return fmt.Errorf("user not found")
	})
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	return ur.repository.Transaction(ctx, func(ctx context.Context) error {
		ur.repository.mutex.Lock()
		defer ur.repository.mutex.Unlock()

		for i, user := range ur.repository.data.Users {
			if user.ID == id {
				ur.repository.data.Users = append(ur.repository.data.Users[:i], ur.repository.data.Users[i+1:]...)
				return nil
			}
		}

		return fmt.Errorf("user not found")
	})
}

func (ur *UserRepository) ListUsers(ctx context.Context) ([]User, error) {
	ur.repository.mutex.RLock()
	defer ur.repository.mutex.RUnlock()

	return ur.repository.data.Users, nil
}
