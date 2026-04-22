package testutil

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"

	"salaryAdvance/internal/infrastructure"
)

func OpenTestPostgresStore(t *testing.T) *infrastructure.PostgresStore {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		raw := os.Getenv("DATABASE_URL")
		if raw != "" {
			dsn = raw
		} else {
			host := getenv("DB_HOST", "localhost")
			port := getenv("DB_PORT", "5432")
			user := getenv("POSTGRES_USER", "loan_user")
			password := getenv("POSTGRES_PASSWORD", "12345678")
			database := getenv("POSTGRES_DB", "loan_service")
			sslmode := getenv("DB_SSLMODE", "disable")
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, database, sslmode)
		}
	}

	store, err := infrastructure.NewPostgresStore(dsn)
	if err != nil {
		t.Skipf("skipping postgres-backed test (db unavailable): %v", err)
	}

	t.Cleanup(func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Errorf("close test postgres store: %v", closeErr)
		}
	})

	return store
}

func TruncateTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec("TRUNCATE TABLE verified_customers, users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
