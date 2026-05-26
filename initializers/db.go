package initializers

import (
	"fmt"
	"os"

	"github.com/TeamPemweb/backkoshub/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		fmt.Println("failed to connect to database")
	}
}

func SyncDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Billing{},
		&models.Complaint{},
		&models.Kamar{},
		&models.TipeKamar{},
		&models.ProfilPemilik{},
		&models.ProfilPenghuni{},)
	if err != nil {
		fmt.Println("failed to migrate database")
	}
}