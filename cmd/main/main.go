package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/iamananya/Ginco-mission-2/pkg/middlewares"
	"github.com/iamananya/Ginco-mission-2/pkg/routes"
)

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret-key"))
	router.Use(sessions.Sessions("session-name", store))
	// CORS ISSUE
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "access-control-allow-headers", "access-control-allow-methods", "access-control-allow-origin", "Session-Id"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Apply authentication middleware to all routes except login and register
	router.Use(middlewares.AuthMiddleware())
	routes.RegisterTicketRoutes(router)

	log.Fatal(router.Run(":9010"))
}
