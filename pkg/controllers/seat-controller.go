package controllers

import (
	"net/http"
	"strconv"

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

func GetSeats(c *gin.Context) {
	seats := models.GetAllSeats()
	c.JSON(http.StatusOK, seats)
}
