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

	kamarIDStr := c.Query("kamar_id")
	penghuniIDStr := c.Query("penghuni_id")

	var kamarID, penghuniID uint64
	var err error

	if kamarIDStr != "" {
		kamarID, err = strconv.ParseUint(kamarIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format kamar_id tidak valid"})
			return
		}
	}

	if penghuniIDStr != "" {
		penghuniID, err = strconv.ParseUint(penghuniIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format penghuni_id tidak valid"})
			return
		}
	}

	billings, err := services.GetAllBillingsByPemilik(pemilikID, uint(kamarID), uint(penghuniID))
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

func UpdateBillingNominal(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	idStr := c.Param("id")
	billingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tagihan tidak valid"})
		return
	}

	var input struct {
		Nominal float64 `json:"nominal" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nominal wajib diisi dan harus lebih besar dari 0"})
		return
	}

	if err := services.UpdateBillingNominal(uint(billingID), pemilikID, input.Nominal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nominal tagihan berhasil diperbarui"})
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

	if err := services.PayBilling(uint(billingID), residentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sistem berhasil mencatat pengalihan konfirmasi ke WhatsApp Owner"})
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