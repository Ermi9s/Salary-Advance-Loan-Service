package services

import (
	"testing"
	"time"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/testutil"
)

func TestAuthRegisterAndLogin(t *testing.T) {
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	repo := repository.NewPostgresRepository(store.DB)
	svc := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:           "test-secret",
		AccessTokenTTL:      time.Hour,
		MaxLoginAttempts:    3,
		LoginWindowDuration: time.Minute,
	})

	err := svc.Register(entity.User{Username: "uploader1", PasswordHash: "Password@123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	token, err := svc.Login("uploader1", "Password@123", "127.0.0.1:uploader1")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	if _, err := svc.ValidateToken(token); err != nil {
		t.Fatalf("token validation failed: %v", err)
	}
}

func TestLoginRateLimit(t *testing.T) {
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	repo := repository.NewPostgresRepository(store.DB)
	svc := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:           "test-secret",
		AccessTokenTTL:      time.Hour,
		MaxLoginAttempts:    2,
		LoginWindowDuration: time.Minute,
	})

	err := svc.Register(entity.User{Username: "user2", PasswordHash: "Password@123"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, _ = svc.Login("user2", "wrong", "ip:user2")
	_, _ = svc.Login("user2", "wrong", "ip:user2")

	_, err = svc.Login("user2", "Password@123", "ip:user2")
	if err == nil {
		t.Fatalf("expected rate limited error")
	}
}
