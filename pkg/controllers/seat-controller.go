package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

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

// Define a lock manager to synchronize seat operations
var seatLocks sync.Map

func CreateSeat(c *gin.Context) {
	var seat models.Seat
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()

	// Start a new transaction
	tx := db.Begin()
	// Acquire a global lock to prevent phantom leads
	globalLock.Lock()
	defer globalLock.Unlock()
	existingSeat := models.GetSeatByNumberAndShowID(tx, seat.SeatNumber, seat.ShowID)
	if existingSeat != nil {
		// Acquire the lock for the existing seat
		lock := getSeatLock(existingSeat.ID)
		lock.Lock()

		// Check if the seat is already booked
		if existingSeat.IsBooked {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Seat is already booked"})
			return
		}

		// Perform optimistic locking by checking the version
		if existingSeat.Version != seat.Version {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Seat has been modified by another user"})
			return
		}

		// Update the existing seat instead of creating a new one
		existingSeat.IsBooked = true
		existingSeat.Version++ // Increment the version
		if err := tx.Save(existingSeat).Error; err != nil {
			lock.Unlock()
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat"})
			return
		}

		tx.Commit()
		lock.Unlock()
		c.JSON(http.StatusOK, existingSeat)
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

// Get or create a lock for the given seat ID
func getSeatLock(seatID uint) *sync.Mutex {
	lock, _ := seatLocks.LoadOrStore(seatID, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

// Global lock to prevent phantom leads
var globalLock sync.Mutex

func DeleteSeat(c *gin.Context) {
	seatNumber := c.Query("seat_number")
	showID := c.Query("show_id")

	db := config.GetDB()

	// Start a new transaction
	tx := db.Begin()
	showIDUint, err := strconv.ParseUint(showID, 10, 64)
	if err != nil {
		// Handle the error, e.g., return an error response or log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showID"})
		return
	}

	seat := models.GetSeatByNumberAndShowID(tx, seatNumber, uint(showIDUint))

	if seat == nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
		return
	}

	if err := tx.Unscoped().Delete(seat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete seat"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Seat deleted"})
}

func UpdateSeat(c *gin.Context) {
	seatNumber := c.Query("seat_number")
	showID := c.Query("show_id")

	db := config.GetDB()

	// Start a new transaction
	tx := db.Begin()
	showIDUint, err := strconv.ParseUint(showID, 10, 64)
	if err != nil {
		// Handle the error, e.g., return an error response or log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showID"})
		return
	}

	seat := models.GetSeatByNumberAndShowID(tx, seatNumber, uint(showIDUint))
	if seat == nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
		return
	}

	seat.IsBooked = true

	if err := tx.Save(seat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat status"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, seat)
}
func BenchmarkSeatBooking(b *testing.B) {
	// Set up the Gin router
	router := gin.Default()
	router.POST("/book-seat", func(c *gin.Context) {
		var seat models.Seat
		fmt.Println("apple start")
		if err := c.ShouldBindJSON(&seat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("apple end")

		db := config.GetDB()
		tx := db.Begin()
		tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

		existingSeat := models.GetSeatByNumberAndShowID(tx, seat.SeatNumber, seat.ShowID)
		if existingSeat != nil && existingSeat.IsBooked {
			tx.Rollback()
			fmt.Println("Seat is already booked")
			c.JSON(http.StatusConflict, gin.H{"error": "Seat is already booked"})
			return
		}

		seat.IsBooked = true
		if err := models.CreateSeat(tx, &seat); err != nil {
			tx.Rollback()
			fmt.Println("Failed to create seat")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seat"})
			return
		}

		tx.Commit()

		c.JSON(http.StatusCreated, seat)
	})

	// Create a new HTTP request for seat booking
	payload := []byte(`{"show_id": 1,
    "seat_number": "C-3",
    "seat_type_id": 1,
    "is_booked": true,
    "user_id":6}`)
	fmt.Println("Payload:", string(payload))

	req, _ := http.NewRequest(http.MethodPost, "/book-seat", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Perform the benchmark test
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Create a test HTTP recorder
			w := httptest.NewRecorder()

			// Serve the HTTP request and record the response
			router.ServeHTTP(w, req)
		}
	})
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
