package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware to check if the user is logged in
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow access to "/login" and "/register" URLs without user ID check
		if c.Request.URL.Path == "/login" || (c.Request.URL.Path == "/register" && c.Request.Method == "POST") {
			c.Next()
			return
		}
		session := sessions.Default(c)
		userID := session.Get("userID")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// User is logged in, continue to the next handler
		c.Next()
	}
}
