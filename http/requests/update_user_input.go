package requests

type UpdateUserInput struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	RememberToken string `json:"remember_token"`
	FCMToken      string `json:"fcm_token"`
}
