package http

import (
	"net/http"
	"salaryAdvance/internal/services"

	"github.com/gin-gonic/gin"
)

type ValidationHandler struct {
	ValidationService *services.ValidationService
}

func (h *ValidationHandler) DataValidationHandler(c *gin.Context) {
	report, err := h.ValidationService.ValidateData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":           "data validated successfully",
		"validation_report": report,
	})
}

func (h *ValidationHandler) GetVerifiedCustomersHandler(c *gin.Context) {
	customers, err := h.ValidationService.GetVerifiedCustomers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":              len(customers),
		"verified_customers": customers,
	})
}
