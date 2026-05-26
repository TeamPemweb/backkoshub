package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/TeamPemweb/backkoshub/initializers"
	"github.com/TeamPemweb/backkoshub/models"
)

func RegisterUser(email, password string) error {
	var existingUser models.User

	err := initializers.DB.
		Where("email = ?", email).
		First(&existingUser).Error

	if err == nil {
		return errors.New("email sudah digunakan")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(hashed),
		Role:     "", 
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AuthenticateUser(email, password string) (string, error) {
	var user models.User

	result := initializers.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("email atau password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,       
		"role": user.Role,     
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type ChooseRoleInput struct {
	Role string `json:"role" binding:"required"`
}

func UpdateUserRole(userID uint, role string) error {
	if role != "pemilik" && role != "penghuni" {
		return errors.New("role harus berupa 'pemilik' atau 'penghuni'")
	}

	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.Role != "" {
		return errors.New("user sudah memiliki role dan tidak dapat diubah")
	}

	user.Role = role
	if err := initializers.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}