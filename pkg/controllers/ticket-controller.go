package controllers

import (
	"net/http"

	"github.com/iamananya/Ginco-mission-2/pkg/models"

	"github.com/gin-gonic/gin"
)

// CreateTicketPrice creates a new ticket price
func CreateTicketPrice(c *gin.Context) {
	ticketPrice := &models.TicketPrice{}
	err := c.ShouldBindJSON(ticketPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tp := models.CreateTicketPrice(ticketPrice)

	c.JSON(http.StatusOK, tp)
}

// GetTicketPrices retrieves all ticket prices
func GetTicketPrices(c *gin.Context) {
	ticketPrices := models.GetAllTicketPrices()
	c.JSON(http.StatusOK, ticketPrices)
}
