package models

import (
	"time"

	"github.com/OdairPianta/gryphon_api/services/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	DefaultModel
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name               string         `gorm:"size:255;notNull" json:"name"`
	Email              string         `gorm:"size:255;notNull;unique" json:"email"`
	EmailVerifiedAt    time.Time      `json:"email_verified_at"`
	Password           string         `gorm:"size:255;notNull" json:"password"`
	RememberToken      string         `gorm:"size:255" json:"remember_token"`
	FCMToken           string         `gorm:"size:255" json:"fcm_token"`
	Token              string         `gorm:"size:255" json:"token"`
	ResetPasswordToken string         `gorm:"size:255" json:"reset_password_token"`
	AwsAccessKeyId     string         `gorm:"size:255" json:"aws_access_key_id"`
	AwsSecretAccessKey string         `gorm:"size:255" json:"aws_secret_access_key"`
	AwsRegion          string         `gorm:"size:255" json:"aws_region"`
	AwsBucket          string         `gorm:"size:255" json:"aws_bucket"`
}

func LoginCheck(email string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserByToken(context *gin.Context) (User, error) {
	var userId, tokenError = token.ExtractTokenID(context)
	if tokenError != nil {
		return User{}, tokenError
	}

	var user User
	err := DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}
