package services

import (
	"context"

	"salaryAdvance/internal/entity"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type ValidationRepository interface {
	SaveVerifiedCustomers(ctx context.Context, customers []entity.BankCustomer) error
	ListVerifiedCustomers(ctx context.Context) ([]entity.BankCustomer, error)
}
