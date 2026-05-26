package models
import (
	"gorm.io/gorm"
)
type TipeKamar struct {
	gorm.Model
	PemilikID      uint      `gorm:"not null"` // Merujuk ke UserID Pemilik
	NamaTipe       string    `gorm:"type:varchar(50);not null"`
	HargaPerBulan  float64   `gorm:"not null"`
	SiklusBayar    int       `gorm:"not null"` // Dalam hitungan bulan (cth: 1, 3, 6, 12)
}