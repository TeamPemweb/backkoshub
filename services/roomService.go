package services

import (
	"errors"
	"math/rand"
	"time"
	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
)

type RoomInput struct {
	TipeKamarID uint   `json:"tipe_kamar_id" binding:"required"`
	NomorKamar  string `json:"nomor_kamar" binding:"required"`
}

func generateKodeKamar(n int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func CreateRoom(pemilikID uint, input RoomInput) error {
	var countTipe int64
	initializers.DB.Model(&models.TipeKamar{}).Where("id = ? AND pemilik_id = ?", input.TipeKamarID, pemilikID).Count(&countTipe)
	if countTipe == 0 {
		return errors.New("tipe kamar tidak ditemukan atau bukan milik Anda")
	}

	var kodeKamar string
	for {
		kodeKamar = generateKodeKamar(6)
		var countKode int64
		initializers.DB.Model(&models.Kamar{}).Where("kode_kamar = ?", kodeKamar).Count(&countKode)
		if countKode == 0 {
			break
		}
	}

	kamar := models.Kamar{
		TipeKamarID: input.TipeKamarID,
		NomorKamar:  input.NomorKamar,
		KodeKamar:   kodeKamar,
		Status:      "kosong",
	}

	return initializers.DB.Create(&kamar).Error
}

func GetRoomsByPemilik(pemilikID uint) ([]models.Kamar, error) {
	var listKamar []models.Kamar
	err := initializers.DB.Preload("TipeKamar").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("tipe_kamars.pemilik_id = ?", pemilikID).
		Find(&listKamar).Error
	return listKamar, err
}

func UpdateRoom(id uint, pemilikID uint, input RoomInput) error {
	var kamar models.Kamar
	err := initializers.DB.Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("kamars.id = ? AND tipe_kamars.pemilik_id = ?", id, pemilikID).
		First(&kamar).Error
	if err != nil {
		return errors.New("kamar tidak ditemukan atau bukan milik Anda")
	}

	var countTipe int64
	initializers.DB.Model(&models.TipeKamar{}).Where("id = ? AND pemilik_id = ?", input.TipeKamarID, pemilikID).Count(&countTipe)
	if countTipe == 0 {
		return errors.New("tipe kamar baru tidak valid atau bukan milik Anda")
	}

	kamar.NomorKamar = input.NomorKamar
	kamar.TipeKamarID = input.TipeKamarID

	return initializers.DB.Save(&kamar).Error
}

func DeleteRoom(id uint, pemilikID uint) error {
	var kamar models.Kamar
	err := initializers.DB.Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("kamars.id = ? AND tipe_kamars.pemilik_id = ?", id, pemilikID).
		First(&kamar).Error
	if err != nil {
		return errors.New("kamar tidak ditemukan atau bukan milik Anda")
	}

	if kamar.Status == "terisi" {
		return errors.New("kamar tidak bisa dihapus karena masih ditempati oleh penghuni")
	}

	return initializers.DB.Delete(&kamar).Error
}
func JoinRoom(residentID uint, kodeKamar string) error {
	var kamar models.Kamar
	err := initializers.DB.Where("kode_kamar = ?", kodeKamar).First(&kamar).Error
	if err != nil {
		return errors.New("kode kamar tidak valid atau tidak ditemukan")
	}

	if kamar.Status == "terisi" {
		return errors.New("kamar tersebut sudah ditempati oleh penghuni lain")
	}

	var count int64
	initializers.DB.Model(&models.Kamar{}).Where("penghuni_id = ?", residentID).Count(&count)
	if count > 0 {
		return errors.New("Anda sudah terdaftar di kamar lain. Silakan keluar dari kamar lama terlebih dahulu")
	}

	now := time.Now()
	kamar.Status = "terisi"
	kamar.PenghuniID = &residentID
	kamar.TanggalMasuk = &now 
	return initializers.DB.Save(&kamar).Error
}

func GetMyRoom(residentID uint) (*models.Kamar, error) {
	var kamar models.Kamar
	err := initializers.DB.Preload("TipeKamar").Where("penghuni_id = ?", residentID).First(&kamar).Error
	if err != nil {
		return nil, errors.New("Anda belum bergabung ke dalam kamar mana pun")
	}
	return &kamar, nil
}