package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/entity"
)

func AuthRequired(authService AuthMiddlewareService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")

		if _, err := authService.ValidateToken(token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		role, err := authService.ParseRole(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid role in token"})
			return
		}

		c.Set("user_role", role)
		c.Set("token", token)
		c.Next()
	}
}

func RequireRole(role entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := c.Get("user_role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role missing"})
			return
		}

		userRole, ok := raw.(entity.UserRole)
		if !ok || userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
