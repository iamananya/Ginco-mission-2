package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/iamananya/Ginco-mission-2/pkg/middlewares"
	"github.com/iamananya/Ginco-mission-2/pkg/routes"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret-key"))
	router.Use(sessions.Sessions("session-name", store))
	// Apply authentication middleware to all routes except login and register

	router.Use(middlewares.AuthMiddleware())
	routes.RegisterTicketRoutes(router)

	log.Fatal(http.ListenAndServe("localhost:9010", router))
}
