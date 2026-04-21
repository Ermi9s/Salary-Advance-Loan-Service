package router

import (
	"github.com/gin-gonic/gin"

	"salaryAdvance/internal/entity"
	handler "salaryAdvance/internal/interfaces/http"
	"salaryAdvance/internal/services"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(
	authService *services.AuthService,
	authHandler *handler.AuthHandlers,
	validationHandler *handler.ValidationHandler,
	ratingHandler *handler.CustomerRatingHandler,
) *Router {
	r := &Router{engine: gin.Default()}

	authGroup := r.engine.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register-admin", handler.AuthRequired(authService), handler.RequireRole(entity.Admin), authHandler.RegisterAdmin)
		authGroup.POST("/logout", handler.AuthRequired(authService), authHandler.Logout)
	}

	secured := r.engine.Group("/api", handler.AuthRequired(authService))
	{
		secured.POST("/validate", validationHandler.DataValidationHandler)
		secured.POST("/process", validationHandler.DataValidationHandler)
		secured.GET("/verified-customers", validationHandler.GetVerifiedCustomersHandler)

		adminOnly := secured.Group("", handler.RequireRole(entity.Admin))
		adminOnly.GET("/customer-ratings", ratingHandler.CustomerRatingHandler)
	}

	return r
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
