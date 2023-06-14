package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iamananya/Ginco-mission-2/pkg/models"
)

var NewUser models.User

// Authenticate user using x-token in headers to get user details-----
func GetUsers(w http.ResponseWriter, r *http.Request) {
	newUsers := models.GetAllUsers()
	res, _ := json.Marshal(newUsers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(requestBody, user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//  Empty Username case has been handled here

	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Print(user.Name, user.Password)
	u := user.CreateUser()

	// Marshal the user object into JSON
	res, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

// GetMovies retrieves all movies
func GetMovies(w http.ResponseWriter, r *http.Request) {
	movies := models.GetAllMovies()
	res, _ := json.Marshal(movies)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetSeatTypes retrieves all seat types
func GetSeatTypes(w http.ResponseWriter, r *http.Request) {
	seatTypes := models.GetAllSeatTypes()
	res, _ := json.Marshal(seatTypes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// CreateTicketPrice creates a new ticket price
func CreateTicketPrice(w http.ResponseWriter, r *http.Request) {
	ticketPrice := &models.TicketPrice{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(requestBody, ticketPrice)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tp := models.CreateTicketPrice(ticketPrice)

	res, err := json.Marshal(tp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetTicketPrices retrieves all ticket prices
func GetTicketPrices(w http.ResponseWriter, r *http.Request) {
	ticketPrices := models.GetAllTicketPrices()
	res, _ := json.Marshal(ticketPrices)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetShows retrieves all shows
func GetShows(w http.ResponseWriter, r *http.Request) {
	shows := models.GetAllShows()
	res, _ := json.Marshal(shows)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetSeats retrieves all seats
func GetSeats(w http.ResponseWriter, r *http.Request) {
	seats := models.GetAllSeats()
	res, _ := json.Marshal(seats)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// CreateBooking creates a new booking
func CreateBooking(w http.ResponseWriter, r *http.Request) {
	booking := &models.Booking{}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(requestBody, booking)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b := models.CreateBookingInitiation(booking)
	if b == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(b)
	fmt.Println("controller", b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// GetBookings retrieves all bookings
func GetBookings(w http.ResponseWriter, r *http.Request) {
	bookings := models.GetAllBookings()
	res, _ := json.Marshal(bookings)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
