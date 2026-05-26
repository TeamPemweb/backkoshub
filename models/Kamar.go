package models

import (
	"gorm.io/gorm"
	"time"
)
type Kamar struct {
	gorm.Model
	TipeKamarID uint      `gorm:"not null"`
	TipeKamar   TipeKamar `gorm:"constraint:OnDelete:RESTRICT;"`
	NomorKamar  string    `gorm:"type:varchar(20);not null"`
	KodeKamar   string    `gorm:"type:varchar(10);unique;not null"` 
	Status      string    `gorm:"type:varchar(20);default:'kosong'"` // "kosong" atau "terisi"
	PenghuniID  *uint     `gorm:"default:null"`
	TanggalMasuk *time.Time `gorm:"default:null"`
}