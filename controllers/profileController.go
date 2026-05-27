package controllers

import (
	"net/http"
	"github.com/TeamPemweb/backkoshub/services"
	"github.com/gin-gonic/gin"
)

type ProfileInput struct {
	Nama         string `json:"nama"`
	NomorTelepon string `json:"nomor_telepon" binding:"required"`
	NamaKos      string `json:"nama_kos"`   
	LokasiKos    string `json:"lokasi_kos"`
}

func GetProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: UserID tidak ditemukan"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe data UserID tidak valid"})
		return
	}

	profile, err := services.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data profil: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func SetupProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Sesi tidak valid"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe data UserID tidak valid"})
		return
	}

	var input ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	err := services.SetupUserProfile(userID, input.Nama, input.NomorTelepon, input.NamaKos, input.LokasiKos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data profil: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil dikonfigurasi",
	})
}
func UpdateProfile(c *gin.Context) {
	userIDInterface, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Sesi tidak valid"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tipe data UserID tidak valid"})
		return
	}

	var input ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	err := services.UpdateUserProfile(userID, input.Nama, input.NomorTelepon, input.NamaKos, input.LokasiKos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui",
	})
}