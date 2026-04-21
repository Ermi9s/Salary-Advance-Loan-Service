package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"salaryAdvance/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type InMemoryRepository struct {
	mu                sync.RWMutex
	users             map[string]entity.User
	verifiedCustomers map[string]entity.BankCustomer
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users:             make(map[string]entity.User),
		verifiedCustomers: make(map[string]entity.BankCustomer),
	}
}

func (r *InMemoryRepository) CreateUser(_ context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Username]; exists {
		return errors.New("username already exists")
	}
	user.CreatedAt = time.Now().UTC()
	r.users[user.Username] = user
	return nil
}

func (r *InMemoryRepository) GetUserByUsername(_ context.Context, username string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[username]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryRepository) SaveVerifiedCustomers(_ context.Context, customers []entity.BankCustomer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, customer := range customers {
		r.verifiedCustomers[customer.AccountNo] = customer
	}
	return nil
}

func (r *InMemoryRepository) ListVerifiedCustomers(_ context.Context) ([]entity.BankCustomer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]entity.BankCustomer, 0, len(r.verifiedCustomers))
	for _, customer := range r.verifiedCustomers {
		result = append(result, customer)
	}
	return result, nil
}
