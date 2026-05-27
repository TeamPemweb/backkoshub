package services

import (
	"errors"
	"time"

	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
)

type BillingResponseDTO struct {
	ID               uint       `json:"id"`
	NomorKamar       string     `json:"nomor_kamar"`
	NamaPenghuni     string     `json:"nama_penghuni"`
	Nominal          float64    `json:"nominal"`
	SiklusBayar      int        `json:"siklus_bayar"`
	JatuhTempo       time.Time  `json:"jatuh_tempo"`
	TanggalBayar     *time.Time `json:"tanggal_bayar"`
	StatusPembayaran string     `json:"status_pembayaran"`
	BuktiPembayaran  string     `json:"bukti_pembayaran"`
}

func GetAllBillingsByPemilik(pemilikID uint, kamarID uint, penghuniID uint) ([]BillingResponseDTO, error) {
	var billings []BillingResponseDTO

	query := initializers.DB.Model(&models.Billing{}).
		Select("billings.id, kamars.nomor_kamar, tipe_kamars.nama_tipe as nama_tipe_kamar, profil_penghunis.nama as nama_penghuni, billings.nominal, billings.siklus_bayar, billings.jatuh_tempo, billings.tanggal_bayar, billings.status_pembayaran, billings.bukti_pembayaran").
		Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Joins("JOIN profil_penghunis ON profil_penghunis.user_id = billings.penghuni_id").
		Where("tipe_kamars.pemilik_id = ?", pemilikID)

	if kamarID > 0 {
		query = query.Where("billings.kamar_id = ?", kamarID)
	}
	if penghuniID > 0 {
		query = query.Where("billings.penghuni_id = ?", penghuniID)
	}

	err := query.Order("billings.created_at DESC").Find(&billings).Error
	return billings, err
}

func ConfirmBillingAsPaid(billingID uint, pemilikID uint) error {
	var billing models.Billing

	err := initializers.DB.Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("billings.id = ? AND tipe_kamars.pemilik_id = ?", billingID, pemilikID).
		First(&billing).Error

	if err != nil {
		return errors.New("tagihan tidak ditemukan atau bukan milik Anda")
	}

	if billing.StatusPembayaran == "lunas" {
		return errors.New("tagihan sudah berstatus lunas")
	}

	now := time.Now()
	billing.StatusPembayaran = "lunas"
	billing.TanggalBayar = &now

	return initializers.DB.Save(&billing).Error
}
type ResidentBillingDTO struct {
	ID               uint       `json:"id"`
	NomorKamar       string     `json:"nomor_kamar"`
	NamaTipeKamar    string    `json:"nama_tipe_kamar"`
	Nominal          float64    `json:"nominal"`
	SiklusBayar      int        `json:"siklus_bayar"`
	JatuhTempo       time.Time  `json:"jatuh_tempo"`
	TanggalBayar     *time.Time `json:"tanggal_bayar"`
	StatusPembayaran string     `json:"status_pembayaran"`
	BuktiPembayaran  string     `json:"bukti_pembayaran"`
}
func UpdateBillingNominal(billingID uint, pemilikID uint, nominal float64) error {
	var billing models.Billing

	err := initializers.DB.Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("billings.id = ? AND tipe_kamars.pemilik_id = ?", billingID, pemilikID).
		First(&billing).Error

	if err != nil {
		return errors.New("tagihan tidak ditemukan atau bukan milik kos Anda")
	}

	return initializers.DB.Model(&billing).Update("nominal", nominal).Error
}

func GetMyActiveBillings(residentID uint) ([]ResidentBillingDTO, error) {
	var billings []ResidentBillingDTO

	err := initializers.DB.Model(&models.Billing{}).
		Select("billings.id, kamars.nomor_kamar, billings.nominal, billings.siklus_bayar, billings.jatuh_tempo, billings.tanggal_bayar, billings.status_pembayaran, billings.bukti_pembayaran").
		Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Where("billings.penghuni_id = ? AND billings.status_pembayaran IN ?", residentID, []string{"menunggu", "lewat_tenggat"}).
		Order("billings.jatuh_tempo ASC").
		Find(&billings).Error

	return billings, err
}

func PayBilling(billingID uint, residentID uint, buktiURL string) error {
	var billing models.Billing

	err := initializers.DB.Where("id = ? AND penghuni_id = ?", billingID, residentID).First(&billing).Error
	if err != nil {
		return errors.New("tagihan tidak ditemukan atau bukan milik Anda")
	}

	if billing.StatusPembayaran == "lunas" {
		return errors.New("tagihan ini sudah lunas")
	}

	billing.BuktiPembayaran = buktiURL
	return initializers.DB.Save(&billing).Error
}

func GetMyBillingHistory(residentID uint) ([]ResidentBillingDTO, error) {
	var billings []ResidentBillingDTO

	err := initializers.DB.Model(&models.Billing{}).
		Select("billings.id, kamars.nomor_kamar, tipe_kamars.nama_tipe as nama_tipe_kamar, billings.nominal, billings.siklus_bayar, billings.jatuh_tempo, billings.tanggal_bayar, billings.status_pembayaran, billings.bukti_pembayaran").
		Joins("JOIN kamars ON kamars.id = billings.kamar_id").
		Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
		Where("billings.penghuni_id = ? AND billings.status_pembayaran = ?", residentID, "lunas").
		Order("billings.tanggal_bayar DESC").
		Find(&billings).Error

	return billings, err
}