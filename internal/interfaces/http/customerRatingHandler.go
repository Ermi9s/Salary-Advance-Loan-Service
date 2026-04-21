package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/services"
)

type CustomerRatingHandler struct {
	ValidationService *services.ValidationService
}

func (h *CustomerRatingHandler) CustomerRatingHandler(c *gin.Context) {
	ratings, err := h.ValidationService.CalculateCustomerRatings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":            len(ratings),
		"customer_ratings": ratings,
	})
}
