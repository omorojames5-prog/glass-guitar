package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/omorojames5-prog/glass-guitar/internal/middleware"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    // Middleware
    router.Use(middleware.CORS())
    
    // Public routes
    router.GET("/health", func(c *gin.Context) {
c.JSON(200, gin.H{
            "status": "ok",
        })
    })

    // Protected routes
    api := router.Group("/api")
    api.Use(middleware.AuthMiddleware())
    {
        // Add your protected routes here
    }

    return router
}
