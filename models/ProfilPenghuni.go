package models
import (
	"gorm.io/gorm"
)
type ProfilPenghuni struct {
	gorm.Model
	UserID       uint   `gorm:"not null"`
	User         User   `gorm:"constraint:OnDelete:CASCADE;"`
	Nama         string `gorm:"type:varchar(100);not null"`
	NomorTelepon string `gorm:"type:varchar(20);not null"`
}