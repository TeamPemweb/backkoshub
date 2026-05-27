package controllers

import (
	"net/http"
	"strconv"

	"github.com/TeamPemweb/backkoshub/services"
	"github.com/gin-gonic/gin"
)

func GetBillings(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	billings, err := services.GetAllBillingsByPemilik(pemilikID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billings)
}

func ConfirmBillingPaid(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	idStr := c.Param("id")
	billingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := services.ConfirmBillingAsPaid(uint(billingID), pemilikID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran tagihan berhasil dikonfirmasi"})
}