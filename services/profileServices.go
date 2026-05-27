package services

import (
	"errors"
	"time"

	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
	"gorm.io/gorm"
)

type ProfileResponse struct {
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Nama         string    `json:"nama"`
	NomorTelepon string    `json:"nomor_telepon"`
	NamaKos      string    `json:"nama_kos,omitempty"`
	LokasiKos    string    `json:"lokasi_kos,omitempty"`
	JoinedAt     time.Time `json:"joined_at"`

	Stats struct {
		TotalTagihanMenunggu int   `json:"total_tagihan_menunggu,omitempty"`
		TotalKomplainPending int   `json:"total_komplain_pending,omitempty"`
		LamaMenetapBulan     int   `json:"lama_menetap_bulan,omitempty"`
		TotalKamarTerisi     int64 `json:"total_kamar_terisi,omitempty"`
	} `json:"stats"`
}

func GetUserProfile(userID uint) (*ProfileResponse, error) {
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	profile := &ProfileResponse{
		Email:    user.Email,
		Role:     user.Role,
		JoinedAt: user.CreatedAt,
	}

	if user.Role == "pemilik" {
		var profilPemilik models.ProfilPemilik
		if err := initializers.DB.Where("user_id = ?", userID).First(&profilPemilik).Error; err == nil {
			profile.Nama = profilPemilik.NamaKos
			profile.NamaKos = profilPemilik.NamaKos
			profile.LokasiKos = profilPemilik.LokasiKos
			profile.NomorTelepon = profilPemilik.NomorTelepon
		}

		var kamarTerisi int64
		initializers.DB.Model(&models.Kamar{}).
			Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
			Where("tipe_kamars.pemilik_id = ? AND kamars.status = ?", userID, "terisi").
			Count(&kamarTerisi)
		profile.Stats.TotalKamarTerisi = kamarTerisi

		var komplainCount int64
		initializers.DB.Model(&models.Complaint{}).
			Joins("JOIN kamars ON kamars.id = complaints.kamar_id").
			Joins("JOIN tipe_kamars ON tipe_kamars.id = kamars.tipe_kamar_id").
			Where("tipe_kamars.pemilik_id = ? AND complaints.status_keluhan = ?", userID, "pending").
			Count(&komplainCount)
		profile.Stats.TotalKomplainPending = int(komplainCount)

	} else if user.Role == "penghuni" {
		var profilPenghuni models.ProfilPenghuni
		if err := initializers.DB.Where("user_id = ?", userID).First(&profilPenghuni).Error; err == nil {
			profile.Nama = profilPenghuni.Nama
			profile.NomorTelepon = profilPenghuni.NomorTelepon
		}

		var tagihanCount int64
		initializers.DB.Model(&models.Billing{}).Where("penghuni_id = ? AND status_pembayaran IN ?", userID, []string{"menunggu", "lewat_tenggat"}).Count(&tagihanCount)
		profile.Stats.TotalTagihanMenunggu = int(tagihanCount)

		var komplainCount int64
		initializers.DB.Model(&models.Complaint{}).
			Joins("JOIN kamars ON kamars.id = complaints.kamar_id").
			Where("kamars.penghuni_id = ? AND complaints.status_keluhan = ?", userID, "pending").
			Count(&komplainCount)
		profile.Stats.TotalKomplainPending = int(komplainCount)

		var kamar models.Kamar
		if err := initializers.DB.Where("penghuni_id = ?", userID).First(&kamar).Error; err == nil && kamar.TanggalMasuk != nil {
			durasi := time.Since(*kamar.TanggalMasuk)
			profile.Stats.LamaMenetapBulan = int(durasi.Hours() / 24 / 30)
		}
	} else {
		return nil, errors.New("role user tidak valid")
	}

	return profile, nil
}

func SetupUserProfile(userID uint, nama, nomorTelepon, namaKos, lokasiKos string) error {
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if user.Role == "pemilik" {
		var profil models.ProfilPemilik
		err := initializers.DB.Where("user_id = ?", userID).First(&profil).Error

		profil.UserID = userID
		profil.NamaKos = namaKos
		profil.LokasiKos = lokasiKos
		profil.NomorTelepon = nomorTelepon

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return initializers.DB.Create(&profil).Error
			}
			return err
		}
		return initializers.DB.Save(&profil).Error

	} else if user.Role == "penghuni" {
		var profil models.ProfilPenghuni
		err := initializers.DB.Where("user_id = ?", userID).First(&profil).Error

		profil.UserID = userID
		profil.Nama = nama
		profil.NomorTelepon = nomorTelepon

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return initializers.DB.Create(&profil).Error
			}
			return err
		}
		return initializers.DB.Save(&profil).Error
	}

	return errors.New("gagal setup profil: role tidak dikenali")
}

func UpdateUserProfile(userID uint, nama, nomorTelepon, namaKos, lokasiKos string) error {
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if user.Role == "pemilik" {
		var profil models.ProfilPemilik
		if err := initializers.DB.Where("user_id = ?", userID).First(&profil).Error; err != nil {
			return errors.New("profil belum dikonfigurasi, silakan lakukan setup profil terlebih dahulu")
		}

		profil.NamaKos = namaKos
		profil.LokasiKos = lokasiKos
		profil.NomorTelepon = nomorTelepon

		return initializers.DB.Save(&profil).Error

	} else if user.Role == "penghuni" {
		var profil models.ProfilPenghuni
		if err := initializers.DB.Where("user_id = ?", userID).First(&profil).Error; err != nil {
			return errors.New("profil belum dikonfigurasi, silakan lakukan setup profil terlebih dahulu")
		}

		profil.Nama = nama
		profil.NomorTelepon = nomorTelepon

		return initializers.DB.Save(&profil).Error
	}

	return errors.New("gagal memperbarui profil: role tidak dikenali")
}