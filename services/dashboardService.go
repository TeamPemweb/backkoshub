package services

import (
	"time"

	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
)

type DashboardStatsResponse struct {
	TotalTunggakan      float64 `json:"total_tunggakan"`
	TotalPenghuniAktif  int64   `json:"total_penghuni_aktif"`
	TotalKomplainPending int64   `json:"total_komplain_pending"`
}

type UnpaidResidentDTO struct {
	BillingID        uint      `json:"billing_id"`
	NomorKamar       string    `json:"nomor_kamar"`
	NamaPenghuni     string    `json:"nama_penghuni"`
	Nominal          float64   `json:"nominal"`
	JatuhTempo       time.Time `json:"jatuh_tempo"`
	StatusPembayaran string    `json:"status_pembayaran"`
}

type ActiveResidentDTO struct {
	KamarID          uint       `json:"kamar_id"`
	NomorKamar       string     `json:"nomor_kamar"`
	NamaTipe         string     `json:"nama_tipe"`
	NamaPenghuni     string     `json:"nama_penghuni"`
	NomorTelepon     string     `json:"nomor_telepon"`
	TanggalMasuk     *time.Time `json:"tanggal_masuk"`
	LamaMenetapBulan int        `json:"lama_menetap_bulan"`
}

func GetDashboardStats(pemilikID uint) (*DashboardStatsResponse, error) {
	var stats DashboardStatsResponse

	initializers.DB.Model(&models.Billing{}).
		Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("tipe_kamars.pemilik_id = ? AND billings.status_pembayaran IN ?", pemilikID, []string{"menunggu", "lewat_tenggat"}).
		Select("COALESCE(SUM(billings.nominal), 0)").
		Scan(&stats.TotalTunggakan)

	initializers.DB.Model(&models.Kamar{}).
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("tipe_kamars.pemilik_id = ? AND kamars.status = ?", pemilikID, "terisi").
		Count(&stats.TotalPenghuniAktif)

	initializers.DB.Model(&models.Complaint{}).
		Joins("JOIN kamars ON kamars.id = complaints.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("tipe_kamars.pemilik_id = ? AND complaints.status_keluhan = ?", pemilikID, "pending").
		Count(&stats.TotalKomplainPending)

	return &stats, nil
}

func GetUnpaidResidents(pemilikID uint) ([]UnpaidResidentDTO, error) {
	var unpaid []UnpaidResidentDTO

	err := initializers.DB.Model(&models.Billing{}).
		Select("billings.id as billing_id, kamars.nomor_kamar, profil_penghunis.nama as nama_penghuni, billings.nominal, billings.jatuh_tempo, billings.status_pembayaran").
		Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Joins("JOIN profil_penghunis ON profil_penghunis.user_id = billings.penghuni_id").
		Where("tipe_kamars.pemilik_id = ? AND billings.status_pembayaran IN ?", pemilikID, []string{"menunggu", "lewat_tenggat"}).
		Order("billings.jatuh_tempo ASC").
		Find(&unpaid).Error

	return unpaid, err
}

func GetActiveResidents(pemilikID uint) ([]ActiveResidentDTO, error) {
	var residents []ActiveResidentDTO

	err := initializers.DB.Model(&models.Kamar{}).
		Select("kamars.id as kamar_id, kamars.nomor_kamar, tipe_kamars.nama_tipe, profil_penghunis.nama as nama_penghuni, profil_penghunis.nomor_telepon, kamars.tanggal_masuk").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Joins("JOIN profil_penghunis ON profil_penghunis.user_id = kamars.penghuni_id").
		Where("tipe_kamars.pemilik_id = ? AND kamars.status = ?", pemilikID, "terisi").
		Find(&residents).Error

	if err != nil {
		return nil, err
	}

	for i := range residents {
		if residents[i].TanggalMasuk != nil {
			durasi := time.Since(*residents[i].TanggalMasuk)
			residents[i].LamaMenetapBulan = int(durasi.Hours() / 24 / 30)
		}
	}

	return residents, nil
}