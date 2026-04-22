package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/interfaces/dto"
	"salaryAdvance/internal/repository"
	"salaryAdvance/internal/services"
	"salaryAdvance/internal/testutil"
)

func TestCustomerRatingHandlerReturnsInternalServerErrorWhenRatingServiceMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := testutil.OpenTestPostgresStore(t)
	testutil.TruncateTables(t, store.DB)
	svc := &services.ValidationService{
		Repository:    repository.NewPostgresRepository(store.DB),
		DTOMethodes:   &dto.DTOMethodes{},
		RatingService: nil,
	}
	h := &CustomerRatingHandler{ValidationService: svc}

	r := gin.New()
	r.GET("/ratings", h.CustomerRatingHandler)

	req := httptest.NewRequest(http.MethodGet, "/ratings", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.Code)
	}
}
