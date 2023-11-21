package requests

type UpdateUserInput struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	RememberToken      string `json:"remember_token"`
	FCMToken           string `json:"fcm_token"`
	AwsAccessKeyId     string `json:"aws_access_key_id"`
	AwsSecretAccessKey string `json:"aws_secret_access_key"`
	AwsRegion          string `json:"aws_region"`
	AwsBucket          string `json:"aws_bucket"`
}
