package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/infrastructure"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user entity.User) error {
	model := infrastructure.UserModel{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
	}

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		if isUniqueViolation(err) {
			return errors.New("username already exists")
		}
		return err
	}

	return nil
}

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var model infrastructure.UserModel
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errors.New("user not found")
		}
		return entity.User{}, err
	}

	return entity.User{
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		Role:         entity.UserRole(model.Role),
		CreatedAt:    model.CreatedAt,
	}, nil
}

func (r *PostgresRepository) SaveVerifiedCustomers(ctx context.Context, customers []entity.BankCustomer) error {
	if len(customers) == 0 {
		return nil
	}

	models := make([]infrastructure.VerifiedCustomerModel, 0, len(customers))
	for _, customer := range customers {
		models = append(models, infrastructure.VerifiedCustomerModel{
			AccountNumber:   customer.AccountNo,
			CustomerName:    customer.CustomerName,
			Mobile:          customer.Mobile,
			BranchName:      customer.BranchName,
			ProductName:     customer.ProductName,
			CustomerID:      customer.CustomerID,
			BranchCode:      customer.BranchCode,
			CustomerBalance: customer.CustomerBalance,
		})
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "account_number"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"customer_name",
				"mobile",
				"branch_name",
				"product_name",
				"customer_id",
				"branch_code",
				"customer_balance",
				"updated_at",
			}),
		}).
		Create(&models).Error
}

func (r *PostgresRepository) ListVerifiedCustomers(ctx context.Context) ([]entity.BankCustomer, error) {
	models := make([]infrastructure.VerifiedCustomerModel, 0)
	err := r.db.WithContext(ctx).Order("id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	customers := make([]entity.BankCustomer, 0, len(models))
	for _, model := range models {
		customers = append(customers, entity.BankCustomer{
			AccountNo:       model.AccountNumber,
			CustomerName:    model.CustomerName,
			Mobile:          model.Mobile,
			BranchName:      model.BranchName,
			ProductName:     model.ProductName,
			CustomerID:      model.CustomerID,
			BranchCode:      model.BranchCode,
			CustomerBalance: model.CustomerBalance,
		})
	}

	return customers, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
