package routes

import (
	"github.com/gorilla/mux"
	"github.com/iamananya/Ginco-mission-2/pkg/controllers"
)

var RegisterTicketRoutes = func(router *mux.Router) {
	router.HandleFunc("/user/", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/user/", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/movies", controllers.GetMovies).Methods("GET")
	router.HandleFunc("/seat-types", controllers.GetSeatTypes).Methods("GET")
	router.HandleFunc("/ticket-prices", controllers.GetTicketPrices).Methods("GET")
	router.HandleFunc("/shows", controllers.GetShows).Methods("GET")
	router.HandleFunc("/seats", controllers.GetSeats).Methods("GET")
	router.HandleFunc("/bookings", controllers.GetBookings).Methods("GET")

	router.HandleFunc("/ticket-prices", controllers.CreateTicketPrice).Methods("POST")
	router.HandleFunc("/bookings", controllers.CreateBooking).Methods("POST")

}
