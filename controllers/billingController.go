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
type PayBillingInput struct {
	BuktiPembayaran string `json:"bukti_pembayaran" binding:"required"`
}

func GetMyBillings(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	residentID := userIDInterface.(uint)

	billings, err := services.GetMyActiveBillings(residentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billings)
}

func PayBilling(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	residentID := userIDInterface.(uint)

	idStr := c.Param("id")
	billingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tagihan tidak valid"})
		return
	}

	var input PayBillingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bukti pembayaran wajib diisi berupa string URL"})
		return
	}

	if err := services.PayBilling(uint(billingID), residentID, input.BuktiPembayaran); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bukti pembayaran berhasil dikirim, menunggu konfirmasi pemilik kos"})
}

func GetMyBillingHistory(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	residentID := userIDInterface.(uint)

	billings, err := services.GetMyBillingHistory(residentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billings)
}