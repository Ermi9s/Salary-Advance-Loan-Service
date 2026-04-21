package repository

import (
	"context"
	"salaryAdvance/internal/entity"
)

type ValidationRepository interface {
	SaveVerifiedCustomers(ctx context.Context, customers []entity.BankCustomer) error
	ListVerifiedCustomers(ctx context.Context) ([]entity.BankCustomer, error)
}
