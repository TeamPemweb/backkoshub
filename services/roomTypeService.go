package services

import (
	"errors"
	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
)

type RoomTypeInput struct {
	NamaTipe      string  `json:"nama_tipe" binding:"required"`
	HargaPerBulan float64 `json:"harga_per_bulan" binding:"required,gt=0"`
	SiklusBayar   int     `json:"siklus_bayar" binding:"required,gt=0"`
}

func CreateRoomType(pemilikID uint, input RoomTypeInput) error {
	tipeKamar := models.TipeKamar{
		PemilikID:     pemilikID,
		NamaTipe:      input.NamaTipe,
		HargaPerBulan: input.HargaPerBulan,
		SiklusBayar:   input.SiklusBayar,
	}

	return initializers.DB.Create(&tipeKamar).Error
}

func GetRoomTypesByPemilik(pemilikID uint) ([]models.TipeKamar, error) {
	var listTipe []models.TipeKamar
	err := initializers.DB.Where("pemilik_id = ?", pemilikID).Find(&listTipe).Error
	return listTipe, err
}

func UpdateRoomType(id uint, pemilikID uint, input RoomTypeInput) error {
	var tipeKamar models.TipeKamar
	if err := initializers.DB.Where("id = ? AND pemilik_id = ?", id, pemilikID).First(&tipeKamar).Error; err != nil {
		return errors.New("tipe kamar tidak ditemukan atau bukan milik Anda")
	}

	tipeKamar.NamaTipe = input.NamaTipe
	tipeKamar.HargaPerBulan = input.HargaPerBulan
	tipeKamar.SiklusBayar = input.SiklusBayar

	if err := initializers.DB.Save(&tipeKamar).Error; err != nil {
		return err
	}
	err := initializers.DB.Model(&models.Billing{}).
		Where("status_pembayaran IN ? AND kamar_id IN (SELECT id FROM kamars WHERE tipe_kamar_id = ?)", []string{"menunggu", "lewat_tenggat"}, id).
		Update("nominal", input.HargaPerBulan).Error
	if err != nil {
		return errors.New("tipe kamar berhasil diperbarui, tetapi gagal menyinkronkan nominal tagihan aktif: " + err.Error())
	}

	return nil
}

func DeleteRoomType(id uint, pemilikID uint) error {
	var tipeKamar models.TipeKamar
	if err := initializers.DB.Where("id = ? AND pemilik_id = ?", id, pemilikID).First(&tipeKamar).Error; err != nil {
		return errors.New("tipe kamar tidak ditemukan atau bukan milik Anda")
	}

	var count int64
	initializers.DB.Model(&models.Kamar{}).Where("tipe_kamar_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("tipe kamar tidak dapat dihapus karena masih digunakan oleh beberapa unit kamar")
	}

	return initializers.DB.Delete(&tipeKamar).Error
}