package infrastructure

import "time"

type UserModel struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"column:username;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Role         string    `gorm:"column:role;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

type VerifiedCustomerModel struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	AccountNumber   string    `gorm:"column:account_number;uniqueIndex;not null"`
	CustomerName    string    `gorm:"column:customer_name;not null"`
	Mobile          string    `gorm:"column:mobile"`
	BranchName      string    `gorm:"column:branch_name"`
	ProductName     string    `gorm:"column:product_name"`
	CustomerID      string    `gorm:"column:customer_id"`
	BranchCode      string    `gorm:"column:branch_code"`
	CustomerBalance float64   `gorm:"column:customer_balance;not null;default:0"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (VerifiedCustomerModel) TableName() string {
	return "verified_customers"
}
