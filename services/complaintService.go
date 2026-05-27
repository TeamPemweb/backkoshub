package services

import (
	"errors"
	"time"

	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
)

type ComplaintResponseDTO struct {
	ID            uint      `json:"id"`
	NomorKamar    string    `json:"nomor_kamar"`
	NamaPenghuni  string    `json:"nama_penghuni"`
	IsiKeluhan    string    `json:"isi_keluhan"`
	StatusKeluhan string    `json:"status_keluhan"`
	CreatedAt     time.Time `json:"created_at"`
}

type UpdateComplaintInput struct {
	StatusKeluhan string `json:"status_keluhan" binding:"required,oneof=pending proses selesai declined"`
}

func GetAllComplaintsByPemilik(pemilikID uint) ([]ComplaintResponseDTO, error) {
	var complaints []ComplaintResponseDTO

	err := initializers.DB.Model(&models.Complaint{}).
		Select("complaints.id, kamars.nomor_kamar, profil_penghunis.nama as nama_penghuni, complaints.isi_keluhan, complaints.status_keluhan, complaints.created_at").
		Joins("JOIN kamars ON kamars.id = complaints.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Joins("JOIN profil_penghunis ON profil_penghunis.user_id = kamars.penghuni_id").
		Where("tipe_kamars.pemilik_id = ?", pemilikID).
		Order("complaints.created_at DESC").
		Find(&complaints).Error

	return complaints, err
}

func UpdateComplaintStatus(complaintID uint, pemilikID uint, status string) error {
	var complaint models.Complaint

	err := initializers.DB.Joins("JOIN kamars ON kamars.id = complaints.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("complaints.id = ? AND tipe_kamars.pemilik_id = ?", complaintID, pemilikID).
		First(&complaint).Error

	if err != nil {
		return errors.New("laporan keluhan tidak ditemukan atau bukan milik kos Anda")
	}

	complaint.StatusKeluhan = status
	return initializers.DB.Save(&complaint).Error
}