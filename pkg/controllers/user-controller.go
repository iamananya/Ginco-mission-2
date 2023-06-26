package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"

	"github.com/iamananya/Ginco-mission-2/pkg/config"
	"github.com/iamananya/Ginco-mission-2/pkg/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// Retrieve the request body
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve the user from the database
	var user models.User
	db := config.GetDB()
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored password
	if user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session ID"})
		return
	}

	// Save the session ID in the session storage
	session := sessions.Default(c)
	fmt.Println("before setting", session.Get("session-id"))
	session.Set("session-id", sessionID)
	err = session.Save()
	fmt.Println("after setting session id", sessionID)

	fmt.Println("after setting", session.Get("session-id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	// Send the session ID in the response headers
	c.Header("Session-ID", sessionID)
	c.Header("Access-Control-Expose-Headers", "session-id")

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
	fmt.Println("Logged In ")
}
func generateSessionID() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a base64 string
	sessionID := base64.URLEncoding.EncodeToString(randomBytes)

	return sessionID, nil
}

func CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	if !VerifyEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email verification failed"})
		return
	}
	createdUser := models.CreateUser(&user)
	c.JSON(http.StatusCreated, createdUser)
}
func VerifyEmail(email string) bool {
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "ananyamahato03@gmail.com"
	senderPassword := "knykuisbyxkmlxmc"

	// Implement email verification logic using an SMTP server
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "ananyamahato03@gmail.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Email Verification for Book My Movie")
	// Generate a verification token
	token := generateVerificationToken()
	// Compose the email body with the verification link
	verificationLink := "http://localhost:3000/login?token=" + token // Replace with your login page URL
	body := "Dear customer, Welcome to Book My Movie. Thanks for choosing us. Excited to be a part of your movie journey.\n.Please click the following link to verify your email:\n\n" + verificationLink
	mailer.SetBody("text/plain", body)

	dialer := gomail.NewDialer(smtpServer, smtpPort, senderEmail, senderPassword)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("Failed to send verification email:", err)
		return false
	}

	return true
}

func generateVerificationToken() string {
	// Generate a unique verification token

	token := uuid.New().String()
	return token
}
func Logout(c *gin.Context) {
	// Delete the session for the logged-out user
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Authenticate user using x-token in headers to get user details-----
func GetUsers(c *gin.Context) {
	users := models.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

// CreateBooking creates a new booking
func CreateBooking(c *gin.Context) {
	booking := &models.Booking{}
	err := c.ShouldBindJSON(booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	b := models.CreateBookingInitiation(booking)
	if b == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, b)
}

// GetBookings retrieves all bookings
func GetBookings(c *gin.Context) {
	bookings := models.GetAllBookings()
	c.JSON(http.StatusOK, bookings)
}
