package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/interfaces/dto"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/services"
	"salaryAdvance/internal/testutil"
)

func TestDataValidationHandlerReturnsInternalServerErrorOnReadFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	svc := &services.ValidationService{
		Repository: repository.NewPostgresRepository(store.DB),
		DTOMethodes: &dto.DTOMethodes{
			CustomersFilePath:      filepath.Join(t.TempDir(), "missing_customers.json"),
			SampleCustomerFilePath: filepath.Join(t.TempDir(), "missing_sample.csv"),
		},
	}
	h := &ValidationHandler{ValidationService: svc}

	r := gin.New()
	r.POST("/validate", h.DataValidationHandler)

	req := httptest.NewRequest(http.MethodPost, "/validate", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.Code)
	}
}

func TestGetVerifiedCustomersHandlerReturnsCustomers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	repo := repository.NewPostgresRepository(store.DB)
	if err := repo.SaveVerifiedCustomers(context.Background(), []entity.BankCustomer{{AccountNo: "123", CustomerName: "Alice"}}); err != nil {
		t.Fatalf("seed verified customer: %v", err)
	}

	svc := &services.ValidationService{Repository: repo}
	h := &ValidationHandler{ValidationService: svc}

	r := gin.New()
	r.GET("/verified", h.GetVerifiedCustomersHandler)

	req := httptest.NewRequest(http.MethodGet, "/verified", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
