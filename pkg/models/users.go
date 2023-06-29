package models

import (
	"fmt"
	"time"

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
	UserID        uint          `json:"user_id" gorm:"foreignKey:UserID"`
	User          User          `json:"user"`
	ShowID        uint          `json:"show_id" gorm:"foreignKey:ShowID"`
	Show          Show          `json:"show"`
	Seats         []Seat        `gorm:"many2many:booking_seats;" json:"seats"`
	TicketDetails TicketDetails `gorm:"foreignKey:BookingID" json:"ticketDetails"`
	PaymentAmount float64       `gorm:"type:decimal(10,2)" json:"payment_amount"`
}
type Transaction struct {
	BookingID    uint      `json:"booking_id"`
	AmountPaid   float64   `json:"amount_paid"`
	SeatsBooked  []Seat    `json:"seats_booked"`
	CreationDate time.Time `json:"creation_date"`
	ShowID       uint      `json:"show_id"`
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
func GetUserByID(userID uint) *User {
	var user User
	if db.First(&user, userID).RecordNotFound() {
		return nil
	}
	return &user
}

func CreateBookingInitiation(booking *Booking, paymentAmount float64) *Booking {
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
	db.Where("user_id = ?", booking.UserID).Find(&seats)

	// Create a slice of strings from the seat numbers
	bookingSeats := make([]Seat, len(seats))
	for i, seat := range seats {
		bookingSeats[i] = seat
	}
	booking.Seats = bookingSeats

	fmt.Println(bookingSeats)
	booking.PaymentAmount = paymentAmount

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
	var booking []Booking
	db.Preload("User").Preload("Show").Find(&booking)

	return booking
}

func GetUserTransactionHistory(userID uint) []Transaction {
	var bookings []Booking
	db.Preload("Seats").Where("user_id = ?", userID).Find(&bookings)

	var transactions []Transaction
	var bookedSeats []Seat
	for _, booking := range bookings {
		var seats []Seat
		db.Model(&booking).Association("Seats").Find(&seats)

		newSeats := filterBookedSeats(seats, bookedSeats)

		if len(newSeats) > 0 {
			transaction := Transaction{
				BookingID:    booking.ID,
				AmountPaid:   booking.PaymentAmount,
				SeatsBooked:  newSeats,
				CreationDate: booking.CreatedAt,
				ShowID:       booking.ShowID,
			}
			transactions = append(transactions, transaction)
			// Append the newly booked seats to the list of booked seats
			bookedSeats = append(bookedSeats, newSeats...)
		}
	}

	return transactions
}

func filterBookedSeats(seats []Seat, bookedSeats []Seat) []Seat {
	var newSeats []Seat

	for _, seat := range seats {
		// Check if the seat is already booked
		isBooked := false
		for _, bookedSeat := range bookedSeats {
			if seat.ID == bookedSeat.ID {
				isBooked = true
				break
			}
		}

		if !isBooked {
			newSeats = append(newSeats, seat)
		}
	}

	return newSeats
}
