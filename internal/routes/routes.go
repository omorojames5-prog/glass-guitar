package routes

import (
"github.com/gin-gonic/gin"
"github.com/omorojames5-prog/glass-guitar/internal/handlers"
"github.com/omorojames5-prog/glass-guitar/internal/middleware"
"github.com/omorojames5-prog/glass-guitar/pkg/database"
)

func SetupRouter() *gin.Engine {
router := gin.Default()

// Middleware
router.Use(middleware.CORS())

// Health check
router.GET("/health", func(c *gin.Context) {
c.JSON(200, gin.H{
"status": "ok",
})
})

// User handlers
userHandler := handlers.NewUserHandler(database.DB)

// Public routes
router.POST("/api/register", userHandler.CreateUser)
router.POST("/api/login", userHandler.Login)

// Protected routes
api := router.Group("/api")
api.Use(middleware.AuthMiddleware())
{
api.GET("/users", userHandler.GetUsers)
api.GET("/users/:id", userHandler.GetUser)
api.PUT("/users/:id", userHandler.UpdateUser)
api.DELETE("/users/:id", userHandler.DeleteUser)
}

return router
}
