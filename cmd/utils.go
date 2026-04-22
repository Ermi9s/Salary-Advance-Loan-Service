package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/infrastructure"
	"salaryAdvance/internal/services"
)

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
	normalizedAdminUser := strings.TrimSpace(strings.ToLower(adminUser))
	if normalizedAdminUser == "" {
		log.Printf("admin seed skipped: username is empty")
		return
	}

	if _, err := authSvc.UserRepo.GetUserByUsername(context.Background(), normalizedAdminUser); err == nil {
		// Admin already exists; skip seeding silently.
		return
	} else if err.Error() != "user not found" {
		log.Printf("admin seed skipped: %v", err)
		return
	}

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

func initializeDatabase() *infrastructure.PostgresStore {
	store, err := infrastructure.NewPostgresStore(buildPostgresDSN())
	if err != nil {
		log.Fatalf("failed to initialize postgres store: %v", err)
	}
	return store
}
