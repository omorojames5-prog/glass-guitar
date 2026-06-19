package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        path := c.Request.URL.Path
        method := c.Request.Method
        
        // Debug log
        println("🔵 AuthMiddleware: Checking path:", path, "method:", method)
        
        // Skip authentication for login and register endpoints
        if path == "/api/login" || path == "/api/register" {
            println("🟢 AuthMiddleware: Skipping auth for:", path)
            c.Next()
            return
        }

        // Skip for OPTIONS requests (CORS preflight)
        if method == "OPTIONS" {
            println("🟢 AuthMiddleware: Skipping OPTIONS request")
            c.Next()
            return
        }

        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            println("🔴 AuthMiddleware: No Authorization header")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            println("🔴 AuthMiddleware: No Bearer token")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        secret := os.Getenv("JWT_SECRET")
        if secret == "" {
            secret = "your-secret-key"
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            println("🔴 AuthMiddleware: Invalid token:", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        println("🟢 AuthMiddleware: Token valid")
        c.Next()
    }
}
