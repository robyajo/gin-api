package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gin-api/internal/delivery/http/handler"
)

// RouterConfig memuat seluruh dependencies handler untuk router
type RouterConfig struct {
	RoleHandler *handler.RoleHandler
	UserHandler *handler.UserHandler
}

// SetupRouter menginisialisasi router Gin dan mendaftarkan seluruh endpoint REST API
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	// CORS basic middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		handler.SuccessResponse(c, http.StatusOK, "Service berjalan normal", gin.H{
			"status": "healthy",
		})
	})

	// API Group v1
	v1 := r.Group("/api/v1")
	{
		// CRUD Roles
		roles := v1.Group("/roles")
		{
			roles.GET("", cfg.RoleHandler.GetRoles)
			roles.POST("", cfg.RoleHandler.CreateRole)
			roles.GET("/:id", cfg.RoleHandler.GetRoleByID)
			roles.PUT("/:id", cfg.RoleHandler.UpdateRole)
			roles.DELETE("/:id", cfg.RoleHandler.DeleteRole)
		}

		// CRUD Users
		users := v1.Group("/users")
		{
			users.GET("", cfg.UserHandler.GetUsers)
			users.POST("", cfg.UserHandler.CreateUser)
			users.GET("/:id", cfg.UserHandler.GetUserByID)
			users.PUT("/:id", cfg.UserHandler.UpdateUser)
			users.DELETE("/:id", cfg.UserHandler.DeleteUser)
		}
	}

	return r
}
