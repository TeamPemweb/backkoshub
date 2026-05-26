package models

import (
	"time"
	"gorm.io/gorm"
)
type Billing struct {
	gorm.Model
	KamarID          uint      `gorm:"not null"`
	Kamar            Kamar     `gorm:"constraint:OnDelete:RESTRICT;"`
	PenghuniID       uint      `gorm:"not null"`
	Nominal          float64   `gorm:"not null"`
	SiklusBayar      int       `gorm:"not null"`
	JatuhTempo       time.Time `gorm:"not null"`
	TanggalBayar     *time.Time`gorm:"default:null"`
	StatusPembayaran string    `gorm:"type:varchar(20);default:'menunggu'"`
	BuktiPembayaran  string    `gorm:"type:text;default:''"`
}