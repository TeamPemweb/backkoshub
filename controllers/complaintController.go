package controllers

import (
	"net/http"
	"strconv"

	"github.com/TeamPemweb/backkoshub/services"
	"github.com/gin-gonic/gin"
)

func GetComplaints(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	complaints, err := services.GetAllComplaintsByPemilik(pemilikID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, complaints)
}

func UpdateComplaintStatus(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	idStr := c.Param("id")
	complaintID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID keluhan tidak valid"})
		return
	}

	var input services.UpdateComplaintInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid. Pilih: pending, proses, selesai, atau declined"})
		return
	}

	if err := services.UpdateComplaintStatus(uint(complaintID), pemilikID, input.StatusKeluhan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status keluhan berhasil diperbarui menjadi '" + input.StatusKeluhan + "'"})
}