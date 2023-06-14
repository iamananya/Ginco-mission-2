package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamananya/Ginco-mission-2/pkg/routes"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterTicketRoutes(r)
	// Apply authentication middleware to all routes
	// authenticatedRouter := middlewares.AuthenticationMiddleware(r.ServeHTTP)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
