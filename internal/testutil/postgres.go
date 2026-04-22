package testutil

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"

	"salaryAdvance/internal/infrastructure"
)

func OpenTestPostgresStore(t testing.TB) *infrastructure.PostgresStore {
	t.Helper()

	candidates := make([]string, 0, 3)
	if dsn := os.Getenv("TEST_DATABASE_URL"); dsn != "" {
		candidates = append(candidates, dsn)
	}

	user := getenv("POSTGRES_USER", "loan_user")
	password := getenv("POSTGRES_PASSWORD", "12345678")
	database := getenv("POSTGRES_DB", "loan_service")
	sslmode := getenv("DB_SSLMODE", "disable")

	candidates = append(candidates,
		fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=%s", user, password, database, sslmode),
		fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=%s", user, password, database, sslmode),
	)

	var lastErr error
	for _, dsn := range candidates {
		store, err := infrastructure.NewPostgresStore(dsn)
		if err == nil {
			t.Cleanup(func() {
				_ = store.Close()
			})
			return store
		}
		lastErr = err
	}

	t.Skipf("skipping postgres-backed test: unable to connect to test db (%v)", lastErr)
	return nil
}

func TruncateTables(t testing.TB, db *gorm.DB) {
	t.Helper()
	if err := db.Exec("TRUNCATE TABLE verified_customers, users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate test tables: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
