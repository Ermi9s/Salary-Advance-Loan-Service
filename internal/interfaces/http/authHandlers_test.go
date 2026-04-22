package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/services"
	"salaryAdvance/internal/testutil"
)

func TestAuthHandlersRegisterLoginLogoutFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	repo := repository.NewPostgresRepository(store.DB)
	authService := services.NewAuthService(repo, services.AuthServiceConfig{
		JWTSecret:           "test-secret",
		AccessTokenTTL:      time.Hour,
		MaxLoginAttempts:    5,
		LoginWindowDuration: time.Minute,
	})

	h := &AuthHandlers{AuthService: authService}
	r := gin.New()
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", AuthRequired(authService), h.Logout)

	registerReq := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"username":"newuser","password":"Password@123"}`))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()
	r.ServeHTTP(registerResp, registerReq)

	if registerResp.Code != http.StatusCreated {
		t.Fatalf("expected register status 201, got %d", registerResp.Code)
	}

	loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"newuser","password":"Password@123"}`))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	r.ServeHTTP(loginResp, loginReq)

	if loginResp.Code != http.StatusOK {
		t.Fatalf("expected login status 200, got %d", loginResp.Code)
	}

	var loginBody map[string]string
	if err := json.Unmarshal(loginResp.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	token := loginBody["token"]
	if token == "" {
		t.Fatalf("expected token in login response")
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+token)
	logoutResp := httptest.NewRecorder()
	r.ServeHTTP(logoutResp, logoutReq)

	if logoutResp.Code != http.StatusOK {
		t.Fatalf("expected logout status 200, got %d", logoutResp.Code)
	}

	if _, err := authService.ValidateToken(token); err == nil {
		t.Fatalf("expected logged out token to be rejected")
	}
}

func TestLogoutWithoutTokenInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	h := &AuthHandlers{AuthService: services.NewAuthService(repository.NewPostgresRepository(store.DB), services.AuthServiceConfig{JWTSecret: "test"})}

	r := gin.New()
	r.POST("/logout", h.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestRegisterAdminSetsAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	repo := repository.NewPostgresRepository(store.DB)
	svc := services.NewAuthService(repo, services.AuthServiceConfig{JWTSecret: "test-secret"})
	h := &AuthHandlers{AuthService: svc}

	r := gin.New()
	r.POST("/auth/register-admin", h.RegisterAdmin)

	req := httptest.NewRequest(http.MethodPost, "/auth/register-admin", strings.NewReader(`{"username":"boss","password":"Admin@1234"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}

	stored, err := repo.GetUserByUsername(context.Background(), "boss")
	if err != nil {
		t.Fatalf("expected admin user to be stored: %v", err)
	}
	if stored.Role != entity.Admin {
		t.Fatalf("expected admin role, got %q", stored.Role)
	}
}
