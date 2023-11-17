package requests

type CreateBase64ImageInput struct {
	Base64    string `json:"base64" binding:"required"`
	Extension string `json:"extension" binding:"required"`
}
