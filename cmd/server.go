package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/infrastructure"
	"salaryAdvance/internal/interfaces/dto"
	handler "salaryAdvance/internal/interfaces/http"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/services"
)

func main() {
	store, err := infrastructure.NewPostgresStore(buildPostgresDSN())
	if err != nil {
		log.Fatalf("failed to initialize postgres store: %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			log.Printf("failed closing postgres store: %v", closeErr)
		}
	}()


	repo := repository.NewPostgresRepository(store.DB)

	authSvc := services.NewAuthService(repo, services.AuthServiceConfig{
		JWTSecret:           getenv("JWT_SECRET", "change-this-in-production"),
		AccessTokenTTL:      2 * time.Hour,
		MaxLoginAttempts:    5,
		LoginWindowDuration: 15 * time.Minute,
	})

	seedAdmin(authSvc)

	dtoMethods := &dto.DTOMethodes{
		CustomersFilePath:      getenv("CUSTOMERS_FILE", "data/customers.json"),
		TransactionFilePath:    getenv("TRANSACTIONS_FILE", "data/transactions.json"),
		SampleCustomerFilePath: getenv("SAMPLE_FILE", "data/sample_customers.csv"),
	}

	validationSvc := &services.ValidationService{
		Repository:    repo,
		DTOMethodes:   dtoMethods,
		RatingService: &services.RatingService{AllowOverdraft: false},
	}

	authHandler := &handler.AuthHandlers{AuthService: authSvc}
	validationHandler := &handler.ValidationHandler{ValidationService: validationSvc}
	ratingHandler := &handler.CustomerRatingHandler{ValidationService: validationSvc}

	r := gin.Default()

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register-admin", handler.AuthRequired(authSvc), handler.RequireRole(entity.Admin), authHandler.RegisterAdmin)
		authGroup.POST("/logout", handler.AuthRequired(authSvc), authHandler.Logout)
	}

	secured := r.Group("/api", handler.AuthRequired(authSvc))
	{
		secured.POST("/validate", validationHandler.DataValidationHandler)
		secured.POST("/process", validationHandler.DataValidationHandler)
		secured.GET("/verified-customers", validationHandler.GetVerifiedCustomersHandler)

		adminOnly := secured.Group("", handler.RequireRole(entity.Admin))
		adminOnly.GET("/customer-ratings", ratingHandler.CustomerRatingHandler)
	}

	port := getenv("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func buildPostgresDSN() string {
	if raw := os.Getenv("DATABASE_URL"); raw != "" {
		return raw
	}

	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("POSTGRES_USER", "loan_user")
	password := getenv("POSTGRES_PASSWORD", "12345678")
	database := getenv("POSTGRES_DB", "loan_service")
	sslmode := getenv("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		database,
		sslmode,
	)
}

func seedAdmin(authSvc *services.AuthService) {
	adminUser := getenv("ADMIN_USERNAME", "admin")
	adminPass := getenv("ADMIN_PASSWORD", "Admin@1234")
	err := authSvc.RegisterAdmin(entity.User{Username: adminUser, PasswordHash: adminPass})
	if err != nil {
		log.Printf("admin seed skipped: %v", err)
	}
}

func getenv(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
