package controllers

import (
	"net/http"

	"github.com/OdairPianta/gryphon_api/http/requests"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/gin-gonic/gin"
)

// @Summary		Add base64 file
// @Description	Add base64 file. Send file_path to save with custom name. Send width and height to resize.
// @Tags			files
// @Accept		json
// @Produce		json
// @Param			file	body		requests.CreateBase64FileInput	true	"Add image"
// @Success		200		{object}	models.File
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/files/base64/create [post]
// @Security Bearer
func CreateBase64File(context *gin.Context) {
	var input requests.CreateBase64FileInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user, error = models.GetUserByToken(context)
	if error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Erro ao obter usuário!"})
		return
	}

	var byteFile []byte
	byteFile, err := models.ConvertBase64IntoByte(input.Base64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var filePath string
	if input.FilePath == "" {
		filePath = models.GenerateRandonFileName(input.Extension)
	} else {
		filePath = input.FilePath
	}

	publicURL, err := models.SaveAsS3(byteFile, input.Extension, user.AwsAccessKeyId, user.AwsSecretAccessKey, user.AwsRegion, user.AwsBucket, filePath)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	model := models.File{
		PublicURL: publicURL,
	}

	context.JSON(http.StatusOK, model)
}
