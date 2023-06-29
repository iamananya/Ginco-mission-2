package controllers

import (
	"net/http"
	"strconv"

	"github.com/iamananya/Ginco-mission-2/pkg/config"
	"github.com/iamananya/Ginco-mission-2/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetMovies retrieves all movies
func GetMovies(c *gin.Context) {
	movies := models.GetAllMovies()
	c.JSON(http.StatusOK, movies)
}

// GetShows retrieves all shows
func GetShows(c *gin.Context) {
	shows := models.GetAllShows()
	c.JSON(http.StatusOK, shows)
}

func GetMovieShows(c *gin.Context) {
	movieIDStr := c.Param("id")
	movieID, err := strconv.ParseUint(movieIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	shows := models.GetShowsByMovieID(uint(movieID))
	c.JSON(http.StatusOK, shows)
}

// GetSeatTypes retrieves all seat types
func GetSeatTypes(c *gin.Context) {
	seatTypes := models.GetAllSeatTypes()
	c.JSON(http.StatusOK, seatTypes)
}
func CreateSeat(c *gin.Context) {
	var seat models.Seat
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := config.GetDB()

	// Start a new transaction
	tx := db.Begin()
	existingSeat := models.GetSeatByNumberAndShowID(tx, seat.SeatNumber, seat.ShowID) // Modify the query to consider both seat number and show ID
	if existingSeat != nil && existingSeat.IsBooked {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Seat is already booked"})
		return
	}

	// Create the seat
	if err := models.CreateSeat(tx, &seat); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seat"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, seat)
}
func GetSeats(c *gin.Context) {
	showID := c.Param("show_id")
	seats := models.GetSeatsByShowID(showID)
	c.JSON(http.StatusOK, seats)
}
func GetMovieByID(c *gin.Context) {
	movieIDStr := c.Param("id")
	movieID, err := strconv.ParseUint(movieIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie := models.GetMovieByID(uint(movieID))
	if movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	shows := models.GetShowsByMovieID(uint(movieID))
	movie.Shows = shows

	c.JSON(http.StatusOK, movie)
}
