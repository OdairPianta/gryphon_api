package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if DB != nil {
		return
	}
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to database!!! error:" + err.Error() + " dsn:" + dsn)
	}
	migrate(database)
	DB = database
}

func migrate(database *gorm.DB) {
	database.Migrator().AutoMigrate(&User{})
	createAdminUser(database)
}

func createAdminUser(database *gorm.DB) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("USER_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to create admin user!!! error:" + err.Error())
	}
	var user User
	database.First(&user, "email = ?", os.Getenv("USER_EMAIL"))
	user.Name = "Admin"
	user.Email = os.Getenv("USER_EMAIL")
	user.EmailVerifiedAt = time.Now()
	user.Password = string(hashedPassword)
	user.FCMToken = ""
	user.AwsAccessKeyId = os.Getenv("USER_AWS_ACCESS_KEY_ID")
	user.AwsSecretAccessKey = os.Getenv("USER_AWS_SECRET_ACCESS_KEY")
	user.AwsRegion = os.Getenv("USER_AWS_REGION")
	user.AwsBucket = os.Getenv("USER_AWS_BUCKET")
	updateErr := database.Save(&user).Error
	if updateErr != nil {
		panic("Failed to create admin user!!! error:" + updateErr.Error())
	}
}
