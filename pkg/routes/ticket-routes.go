package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iamananya/Ginco-mission-2/pkg/controllers"
)

var RegisterTicketRoutes = func(router *gin.Engine) {
	router.POST("/register", controllers.CreateUser)
	router.GET("/user", controllers.GetUsers)
	router.GET("/user/:id", controllers.GetSingleUser)

	router.GET("/movies", controllers.GetMovies)
	router.GET("/seat-types", controllers.GetSeatTypes)
	router.GET("/ticket-prices", controllers.GetTicketPrices)
	router.GET("/shows", controllers.GetShows)
	router.GET("/seats/:show_id", controllers.GetSeats)
	router.GET("/bookings", controllers.GetBookings)
	router.GET("/movies/:id", controllers.GetMovieByID)
	router.POST("/seats", controllers.CreateSeat)
	router.GET("/transaction-history/:userID", controllers.GetUserTransactionHistory)

	router.POST("/ticket-prices", controllers.CreateTicketPrice)
	router.POST("/bookings", controllers.CreateBooking)
	router.POST("/logout", controllers.Logout)
}
