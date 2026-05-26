package models
import (
	"gorm.io/gorm"
)
type ProfilPemilik struct {
	gorm.Model
	UserID       uint   `gorm:"not null"`
	User         User   `gorm:"constraint:OnDelete:CASCADE;"`
	NamaKos      string `gorm:"type:varchar(100);not null"`
	LokasiKos    string `gorm:"type:text;not null"`
	NomorTelepon string `gorm:"type:varchar(20);not null"`
}
