package requests

type CreateUserInput struct {
	Name               string `json:"name" binding:"required"`
	Email              string `json:"email" binding:"required"`
	Password           string `json:"password" binding:"required"`
	RememberToken      string `json:"remember_token"`
	FCMToken           string `json:"fcm_token"`
	AwsAccessKeyId     string `json:"aws_access_key_id" binding:"required"`
	AwsSecretAccessKey string `json:"aws_secret_access_key" binding:"required"`
	AwsRegion          string `json:"aws_region" binding:"required"`
	AwsBucket          string `json:"aws_bucket" binding:"required"`
}
