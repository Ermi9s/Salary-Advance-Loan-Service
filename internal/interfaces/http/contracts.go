package http

import (
	"context"

	"github.com/golang-jwt/jwt/v5"

	"salaryAdvance/internal/entity"
)

type AuthHandlerService interface {
	Register(user entity.User) error
	RegisterAdmin(user entity.User) error
	Login(username, password, sourceKey string) (string, error)
	Logout(token string) error
}

type AuthMiddlewareService interface {
	ValidateToken(tokenString string) (*jwt.RegisteredClaims, error)
	ParseRole(tokenString string) (entity.UserRole, error)
}

type ValidationHandlerService interface {
	ValidateData(ctx context.Context) (entity.ValidationReport, error)
	GetVerifiedCustomers(ctx context.Context) ([]entity.BankCustomer, error)
}

type CustomerRatingService interface {
	CalculateCustomerRatings(ctx context.Context) ([]entity.CustomerRating, error)
}
