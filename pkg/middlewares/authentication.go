package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware to check if the user is logged in
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Next()
			return
		}
		session := sessions.Default(c)
		userID := session.Get("userID")
		if userID == nil {
			// User is not logged in, redirect to the login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// User is logged in, continue to the next handler
		c.Next()
	}
}
