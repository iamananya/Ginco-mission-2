package models

import (
	"fmt"
	"time"

	"github.com/iamananya/Ginco-mission-2/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(30);size:30" json:"name"`
	Password string `gorm:"type:char(60)" json:"password"`
	Email    string `gorm:"" json:"email"`
}
type Movie struct {
	gorm.Model
	Title  string  `gorm:"type:varchar(30)" json:"title"`
	Genre  string  `gorm:"type:varchar(30)" json:"genre"`
	Rating float64 `gorm:"type:decimal(10,2)" json:"rating"`
}
type SeatType struct {
	gorm.Model
	Name         string `gorm:"type:varchar(30);size:30" json:"name"`
	Description  string `gorm:"type:varchar(100)" json:"description"`
	TicketPrices []TicketPrice
}
type TicketPrice struct {
	gorm.Model
	MovieID    uint    `json:"movie_id"`
	SeatTypeID uint    `json:"seat_type_id"`
	Price      float64 `gorm:"type:decimal(10,2)" json:"price"`
	Shows      []Show
}
type Show struct {
	gorm.Model
	MovieID       uint      `json:"movie_id"`
	Showtime      time.Time `json:"showtime"`
	TicketPriceID uint      `json:"ticket_price_id"`
	Seats         []Seat
}
type Seat struct {
	gorm.Model
	ShowID     uint   `json:"show_id"`
	SeatNumber string `gorm:"type:varchar(10)" json:"seat_number"`
	SeatTypeID uint   `json:"seat_type_id"`
	IsBooked   bool   `json:"is_booked"`
}
type Booking struct {
	gorm.Model
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignkey:UserID" json:"-"`
	ShowID uint `json:"show_id"`
	Show   Show `gorm:"foreignkey:ShowID" json:"-"`

	Seats         []Seat        `gorm:"many2many:booking_seats;" json:"seats"`
	TicketDetails TicketDetails `gorm:"foreignKey:BookingID" json:"ticketDetails"`
	PaymentAmount float64       `gorm:"type:decimal(10,2)" json:"payment_amount"`
}

type TicketDetails struct {
	gorm.Model
	BookingID uint `json:"booking_id"`
	User      User `json:"user"`
	Showtime  Show `json:"showtime"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{}, &Movie{}, &SeatType{}, &TicketPrice{}, &Show{}, &Seat{}, &Booking{}, &TicketDetails{})
	db.Model(&User{}).ModifyColumn("password", "char(60)")
	// Define the relationships
	db.Model(&SeatType{}).Related(&TicketPrice{}, "TicketPrices")
	db.Model(&TicketPrice{}).Related(&Show{}, "Shows")
	db.Model(&Show{}).Related(&Seat{}, "Seats")
}

func CreateUser(u *User) *User {
	db.NewRecord(u)
	db.Create(&u)
	return u
}

func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
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
func CreateTicketPrice(ticketPrice *TicketPrice) *TicketPrice {
	db.NewRecord(ticketPrice)
	db.Create(&ticketPrice)
	return ticketPrice
}
func GetAllTicketPrices() []TicketPrice {
	var ticketPrices []TicketPrice
	db.Find(&ticketPrices)
	return ticketPrices
}
func GetAllShows() []Show {
	var shows []Show
	db.Find(&shows)
	return shows
}

func GetAllSeats() []Seat {
	var seats []Seat
	db.Find(&seats)
	return seats
}

func CreateBookingInitiation(booking *Booking) *Booking {
	// Validate the request body
	if booking.UserID == 0 || booking.ShowID == 0 {
		return nil
	}

	// Retrieve user details
	var user User
	db.First(&user, booking.UserID)
	booking.User = user
	fmt.Println("User", booking.User)
	// Retrieve show details
	var show Show
	db.First(&show, booking.ShowID)
	booking.Show = show
	fmt.Println("Show", booking.Show)
	// Get the seat details
	var seats []Seat
	db.Where("show_id = ?", booking.ShowID).Find(&seats)

	// Create a slice of strings from the seat numbers
	bookingSeats := make([]Seat, len(seats))
	for i, seat := range seats {
		bookingSeats[i] = seat
	}
	booking.Seats = bookingSeats

	fmt.Println(bookingSeats)
	// Calculate the payment amount
	var ticketPrice TicketPrice
	db.First(&ticketPrice, show.TicketPriceID)
	booking.PaymentAmount = calculatePaymentAmount(len(bookingSeats), ticketPrice.Price)
	fmt.Println("Paymebt", booking.PaymentAmount)

	// Create the ticket details
	ticketDetails := TicketDetails{
		User:     user,
		Showtime: show,
	}

	// Create a slice of TicketDetails values
	booking.TicketDetails = ticketDetails
	// fmt.Println(booking)
	fmt.Println("*************************")
	fmt.Println("Booking Models:", booking)

	// Create the booking
	db.Save(&booking)

	return booking
}

func GetAllBookings() []Booking {
	var bookings []Booking
	db.Preload("User").Preload("Show").Find(&bookings)

	return bookings
}

func calculatePaymentAmount(numTickets int, ticketPrice float64) float64 {
	return float64(numTickets) * ticketPrice
}
