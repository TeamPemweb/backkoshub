package controllers

import (
	"net/http"

	"github.com/TeamPemweb/backkoshub/services"
	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	stats, err := services.GetDashboardStats(pemilikID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetUnpaidResidents(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	unpaidList, err := services.GetUnpaidResidents(pemilikID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, unpaidList)
}

func GetActiveResidents(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	residentsList, err := services.GetActiveResidents(pemilikID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, residentsList)
}