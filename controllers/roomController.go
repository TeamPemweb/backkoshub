package controllers

import (
	"net/http"
	"strconv"
	"github.com/TeamPemweb/backkoshub/services"
	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe data UserID tidak valid"})
		return
	}

	var input services.RoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	if err := services.CreateRoom(pemilikID, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kamar berhasil ditambahkan"})
}

func GetRooms(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	listKamar, err := services.GetRoomsByPemilik(pemilikID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listKamar)
}

func UpdateRoom(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input services.RoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	if err := services.UpdateRoom(uint(id), pemilikID, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kamar berhasil diperbarui"})
}

func DeleteRoom(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := services.DeleteRoom(uint(id), pemilikID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kamar berhasil dihapus"})
}
func JoinRoom(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	residentID := userIDInterface.(uint)

	var input struct {
		KodeKamar string `json:"kode_kamar" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode kamar wajib diisi"})
		return
	}

	if err := services.JoinRoom(residentID, input.KodeKamar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil bergabung ke dalam kamar"})
}

func GetMyRoom(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	residentID := userIDInterface.(uint)

	kamar, err := services.GetMyRoom(residentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kamar)
}
func EndLease(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi tidak valid"})
		return
	}
	pemilikID := userIDInterface.(uint)

	idStr := c.Param("id")
	kamarID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kamar tidak valid"})
		return
	}

	if err := services.EndLease(uint(kamarID), pemilikID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sewa berhasil diakhiri, status kamar kembali kosong"})
}