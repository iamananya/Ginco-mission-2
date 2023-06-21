package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Movie struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(30)" json:"title"`
	Genre       string  `gorm:"type:varchar(30)" json:"genre"`
	Rating      float64 `gorm:"type:decimal(10,2)" json:"rating"`
	Image       string  `gorm:"type:varchar(255)" json:"image"`
	Description string  `gorm:"type:text" json:"desc"`
	Shows       []Show  `json:"shows"`
}
type Show struct {
	gorm.Model
	MovieID       uint      `json:"movie_id"`
	Showtime      time.Time `json:"showtime"`
	TicketPriceID uint      `json:"ticket_price_id"`
	Seats         []Seat
}
type SeatType struct {
	gorm.Model
	Name         string `gorm:"type:varchar(30);size:30" json:"name"`
	Description  string `gorm:"type:varchar(100)" json:"description"`
	TicketPrices []TicketPrice
}
type Seat struct {
	gorm.Model
	ShowID     uint   `json:"show_id"`
	SeatNumber string `gorm:"type:varchar(10)" json:"seat_number"`
	SeatTypeID uint   `json:"seat_type_id"`
	IsBooked   bool   `json:"is_booked"`
}

func GetAllMovies() []Movie {
	var movies []Movie
	db.Find(&movies)
	return movies
}
func GetAllSeatTypes() []SeatType {
	var seatTypes []SeatType
	db.Find(&seatTypes)
	return seatTypes
}

func GetAllSeats() []Seat {
	var seats []Seat
	db.Find(&seats)
	return seats
}
func GetAllShows() []Show {
	var shows []Show
	db.Find(&shows)
	return shows
}
func GetShowsByMovieID(movieID uint) []Show {
	var shows []Show
	db.Where("movie_id = ?", movieID).Find(&shows)
	return shows
}
