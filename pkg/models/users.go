package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(30);size:30" json:"name"`
	Password string `gorm:"type:varchar(255)" json:"password"`
	Email    string `gorm:"" json:"email"`
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
