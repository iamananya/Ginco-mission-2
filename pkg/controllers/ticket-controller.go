package controllers

import (
	"net/http"

	"github.com/iamananya/Ginco-mission-2/pkg/config"
	"github.com/iamananya/Ginco-mission-2/pkg/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	// Retrieve the request body
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve the user from the database
	var user models.User
	db := config.GetDB()
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create a session for the logged-in user
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func Logout(c *gin.Context) {
	// Delete the session for the logged-out user
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Authenticate user using x-token in headers to get user details-----
func GetUsers(c *gin.Context) {
	users := models.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.Password = string(hashedPassword)
	createdUser := models.CreateUser(&user)
	c.JSON(http.StatusCreated, createdUser)

}

// GetMovies retrieves all movies
func GetMovies(c *gin.Context) {
	movies := models.GetAllMovies()
	c.JSON(http.StatusOK, movies)
}

// GetSeatTypes retrieves all seat types
func GetSeatTypes(c *gin.Context) {
	seatTypes := models.GetAllSeatTypes()
	c.JSON(http.StatusOK, seatTypes)
}

// CreateTicketPrice creates a new ticket price
func CreateTicketPrice(c *gin.Context) {
	ticketPrice := &models.TicketPrice{}
	err := c.ShouldBindJSON(ticketPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tp := models.CreateTicketPrice(ticketPrice)

	c.JSON(http.StatusOK, tp)
}

// GetTicketPrices retrieves all ticket prices
func GetTicketPrices(c *gin.Context) {
	ticketPrices := models.GetAllTicketPrices()
	c.JSON(http.StatusOK, ticketPrices)
}

// GetShows retrieves all shows
func GetShows(c *gin.Context) {
	shows := models.GetAllShows()
	c.JSON(http.StatusOK, shows)
}

// GetSeats retrieves all seats
func GetSeats(c *gin.Context) {
	seats := models.GetAllSeats()
	c.JSON(http.StatusOK, seats)
}

// CreateBooking creates a new booking
func CreateBooking(c *gin.Context) {
	booking := &models.Booking{}
	err := c.ShouldBindJSON(booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	b := models.CreateBookingInitiation(booking)
	if b == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, b)
}

// GetBookings retrieves all bookings
func GetBookings(c *gin.Context) {
	bookings := models.GetAllBookings()
	c.JSON(http.StatusOK, bookings)
}

func SetupSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret"))
	return sessions.Sessions("session", store)
}
