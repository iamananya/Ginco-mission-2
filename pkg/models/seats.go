package models

import (
	"errors"
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
	MovieID        uint      `json:"movie_id"`
	Showtime_slot1 time.Time `json:"showtime_slot_1"`
	Showtime_slot2 time.Time `json:"showtime_slot_2"`
	Showtime_slot3 time.Time `json:"showtime_slot_3"`
	Showtime_slot4 time.Time `json:"showtime_slot_4"`

	TicketPriceID uint `json:"ticket_price_id"`
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
	UserID     uint   `json:"user_id"`
	Version    uint   `json:"version" gorm:"default:1"`
}

func GetAllMovies() []Movie {
	var movies []Movie
	db.Preload("Shows").Find(&movies)
	return movies
}
func GetAllSeatTypes() []SeatType {
	var seatTypes []SeatType
	db.Find(&seatTypes)
	return seatTypes
}

var ErrSeatAlreadyBooked = errors.New("seat is already booked")

func GetSeatByNumber(tx *gorm.DB, seatNumber string) *Seat {
	var seat Seat
	if err := tx.Where("seat_number = ?", seatNumber).First(&seat).Error; err != nil {
		return nil
	}
	return &seat
}
func CreateSeat(tx *gorm.DB, seat *Seat) error {
	if err := tx.Create(seat).Error; err != nil {
		return err
	}
	return nil
}
func GetAllSeats() []Seat {
	var seats []Seat
	db.Find(&seats)
	return seats
}
func GetSeatByNumberAndShowID(tx *gorm.DB, seatNumber string, showID uint) *Seat {
	var seat Seat
	if err := tx.Where("seat_number = ? AND show_id = ?", seatNumber, int(showID)).First(&seat).Error; err != nil {
		return nil
	}
	return &seat
}
func GetSeatsByShowID(showID string) []Seat {
	var seats []Seat
	db.Where("show_id = ?", showID).Find(&seats)
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
func GetMovieByID(movieID uint) *Movie {
	var movie Movie
	if db.First(&movie, movieID).RecordNotFound() {
		return nil
	}
	return &movie
}
