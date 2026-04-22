package main

import (
	"log"
	"time"
	"salaryAdvance/internal/interfaces/dto"
	handler "salaryAdvance/internal/interfaces/http"
	"salaryAdvance/internal/interfaces/router"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/services"
)


func main() {
	store := initializeDatabase()
	repo := repository.NewPostgresRepository(store.DB)
	authSvc := services.NewAuthService(repo, services.AuthServiceConfig{
		JWTSecret:           getenv("JWT_SECRET", "change-this-in-production"),
		AccessTokenTTL:      2 * time.Hour,
		MaxLoginAttempts:    5,
		LoginWindowDuration: 15 * time.Minute,
	})
	
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
	
	seedAdmin(authSvc)

	authHandler := &handler.AuthHandlers{AuthService: authSvc}
	validationHandler := &handler.ValidationHandler{ValidationService: validationSvc}
	ratingHandler := &handler.CustomerRatingHandler{ValidationService: validationSvc}

	r := router.NewRouter(authSvc, authHandler, validationHandler, ratingHandler)

	port := getenv("PORT", "8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

