package models

import (
	"github.com/iamananya/Ginco-mission-2/pkg/config"
	"github.com/jinzhu/gorm"
)

type TicketPrice struct {
	gorm.Model
	MovieID    uint    `json:"movie_id"`
	SeatTypeID uint    `json:"seat_type_id"`
	Price      float64 `gorm:"type:decimal(10,2)" json:"price"`
	Shows      []Show
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
	db.Model(&User{}).ModifyColumn("password", "varchar(255)")
	db.Model(&User{}).ModifyColumn("desc", "text")

	db.Model(&Movie{}).AutoMigrate(&Movie{})
	// Define the relationships
	db.Model(&SeatType{}).Related(&TicketPrice{}, "TicketPrices")
	db.Model(&TicketPrice{}).Related(&Show{}, "Shows")
	db.Model(&Show{}).Related(&Seat{}, "Seats")
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
