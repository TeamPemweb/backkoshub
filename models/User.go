package models

import (
	"gorm.io/gorm"
	"time"
)
type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"type:varchar(20);not null"` // "pemilik" atau "penghuni"
}
type PasswordResetToken struct {
	gorm.Model
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"type:varchar(100);unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}