package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/testutil"
)

func TestIsUniqueViolation(t *testing.T) {
	if !isUniqueViolation(&pgconn.PgError{Code: "23505"}) {
		t.Fatalf("expected unique violation to be detected")
	}

	wrapped := fmt.Errorf("wrapped error: %w", &pgconn.PgError{Code: "23505"})
	if !isUniqueViolation(wrapped) {
		t.Fatalf("expected wrapped unique violation to be detected")
	}

	if isUniqueViolation(&pgconn.PgError{Code: "12345"}) {
		t.Fatalf("did not expect non-unique pg error to match")
	}

	if isUniqueViolation(nil) {
		t.Fatalf("did not expect nil error to match")
	}
}

func TestPostgresRepositoryCreateAndGetUser(t *testing.T) {
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)

	repo := NewPostgresRepository(store.DB)
	ctx := context.Background()

	user := entity.User{Username: "alice", PasswordHash: "hash", Role: entity.Uploader}
	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	stored, err := repo.GetUserByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}

	if stored.Username != "alice" {
		t.Fatalf("unexpected username: %q", stored.Username)
	}
	if stored.Role != entity.Uploader {
		t.Fatalf("unexpected role: %q", stored.Role)
	}
	if stored.CreatedAt.IsZero() {
		t.Fatalf("expected created timestamp to be set")
	}

	if err := repo.CreateUser(ctx, user); err == nil {
		t.Fatalf("expected duplicate username error")
	}
}

func TestPostgresRepositorySaveAndListVerifiedCustomers(t *testing.T) {
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)

	repo := NewPostgresRepository(store.DB)
	ctx := context.Background()

	initial := []entity.BankCustomer{
		{AccountNo: "111", CustomerName: "Alice", CustomerBalance: 100},
		{AccountNo: "222", CustomerName: "Bob", CustomerBalance: 200},
	}
	if err := repo.SaveVerifiedCustomers(ctx, initial); err != nil {
		t.Fatalf("SaveVerifiedCustomers initial failed: %v", err)
	}

	update := []entity.BankCustomer{{AccountNo: "222", CustomerName: "Bobby", CustomerBalance: 250}}
	if err := repo.SaveVerifiedCustomers(ctx, update); err != nil {
		t.Fatalf("SaveVerifiedCustomers upsert failed: %v", err)
	}

	customers, err := repo.ListVerifiedCustomers(ctx)
	if err != nil {
		t.Fatalf("ListVerifiedCustomers failed: %v", err)
	}

	if len(customers) != 2 {
		t.Fatalf("expected 2 customers, got %d", len(customers))
	}

	foundUpdated := false
	for _, c := range customers {
		if c.AccountNo == "222" {
			foundUpdated = true
			if c.CustomerName != "Bobby" || c.CustomerBalance != 250 {
				t.Fatalf("expected updated customer data, got %#v", c)
			}
		}
	}

	if !foundUpdated {
		t.Fatalf("expected account 222 in listed customers")
	}
}
