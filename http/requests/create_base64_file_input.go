package requests

type CreateBase64FileInput struct {
	Base64    string `json:"base64" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	FilePath  string `json:"file_path"`
}
