package models
import (
	"gorm.io/gorm"
)
type Complaint struct {
	gorm.Model
	KamarID        uint      `gorm:"not null"`
	PenghuniID     uint      `gorm:"not null"`
	IsiKeluhan     string    `gorm:"type:text;not null"`
	StatusKeluhan  string    `gorm:"type:varchar(20);default:'pending'"` // "pending", "accepted", "declined"
}