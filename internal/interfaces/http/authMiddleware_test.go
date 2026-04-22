package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/services"
	"salaryAdvance/internal/testutil"
)

func newTestAuthService(t *testing.T) *services.AuthService {
	t.Helper()

	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	repo := repository.NewPostgresRepository(store.DB)
	svc := services.NewAuthService(repo, services.AuthServiceConfig{
		JWTSecret:           "test-secret",
		AccessTokenTTL:      time.Hour,
		MaxLoginAttempts:    3,
		LoginWindowDuration: time.Minute,
	})

	if err := svc.RegisterAdmin(entity.User{Username: "admin", PasswordHash: "Admin@1234"}); err != nil {
		t.Fatalf("RegisterAdmin failed: %v", err)
	}
	if err := svc.Register(entity.User{Username: "uploader", PasswordHash: "Uploader@123"}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	return svc
}

func tokenFor(t *testing.T, svc *services.AuthService, username, password string) string {
	t.Helper()
	token, err := svc.Login(username, password, "ip:"+username)
	if err != nil {
		t.Fatalf("Login failed for %s: %v", username, err)
	}
	return token
}

func TestAuthRequiredRejectsMissingBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newTestAuthService(t)

	r := gin.New()
	r.GET("/protected", AuthRequired(svc), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestAuthRequiredSetsRoleAndTokenInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newTestAuthService(t)
	uploaderToken := tokenFor(t, svc, "uploader", "Uploader@123")

	r := gin.New()
	r.GET("/protected", AuthRequired(svc), func(c *gin.Context) {
		role, _ := c.Get("user_role")
		token, _ := c.Get("token")
		c.JSON(http.StatusOK, gin.H{"role": role, "token": token})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+uploaderToken)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["role"] != string(entity.Uploader) {
		t.Fatalf("expected role %q, got %#v", entity.Uploader, payload["role"])
	}
	if payload["token"] != uploaderToken {
		t.Fatalf("expected token to be set in context")
	}
}

func TestRequireRoleEnforcesAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uploaderRouter := gin.New()
	uploaderRouter.GET("/admin", func(c *gin.Context) {
		c.Set("user_role", entity.Uploader)
		c.Next()
	}, RequireRole(entity.Admin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	uploaderReq := httptest.NewRequest(http.MethodGet, "/admin", nil)
	uploaderResp := httptest.NewRecorder()
	uploaderRouter.ServeHTTP(uploaderResp, uploaderReq)
	if uploaderResp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for uploader, got %d", uploaderResp.Code)
	}

	adminRouter := gin.New()
	adminRouter.GET("/admin", func(c *gin.Context) {
		c.Set("user_role", entity.Admin)
		c.Next()
	}, RequireRole(entity.Admin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	adminReq := httptest.NewRequest(http.MethodGet, "/admin", nil)
	adminResp := httptest.NewRecorder()
	adminRouter.ServeHTTP(adminResp, adminReq)
	if adminResp.Code != http.StatusOK {
		t.Fatalf("expected 200 for admin, got %d", adminResp.Code)
	}
}
