package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware to check if the user is logged in
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" || c.Request.URL.Path == "/movies" || strings.HasPrefix(c.Request.URL.Path, "/movies/") || c.Request.URL.Path == "/verify-email" {
			c.Next()
			return
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		// session := sessions.Default(c)
		// sessionID := session.Get("session-id")
		requestSessionID := c.GetHeader("Session-ID")
		fmt.Println("from request session id", requestSessionID)

		// fmt.Println("session id from session storage", sessionID)
		// if requestSessionID == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "request empty Unauthorized"})
		// 	c.Abort()
		// 	return
		// }
		// if sessionID == nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "empty Unauthorized"})
		// 	c.Abort()
		// 	return
		// }

		// // Compare the session ID from the session storage with the one in the request headers
		// if sessionID != requestSessionID {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "not same Unauthorized"})
		// 	c.Abort()
		// 	return
		// }

		// User is authenticated, continue to the next handler
		c.Next()
	}
}
