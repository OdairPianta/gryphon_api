package controllers

import (
	"net/http"

	"github.com/OdairPianta/gryphon_api/http/requests"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/gin-gonic/gin"
)

// @Summary		Add base64 image
// @Description	Add base64 image. Send file_path to save with custom name. Send Width and Height to resize.
// @Tags			images
// @Accept		json
// @Produce		json
// @Param			image	body		requests.CreateBase64ImageInput	true	"Add image"
// @Success		200		{object}	models.Image
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/images [post]
// @Security Bearer
func CreateBase64Image(context *gin.Context) {
	var input requests.CreateBase64ImageInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user, error = models.GetUserByToken(context)
	if error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Erro ao obter usuário!"})
		return
	}

	image, err := models.ConvertBase64IntoImage(input.Base64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var byteImage []byte
	if input.Width != 0 || input.Height != 0 {
		resizedImage := models.ResizeImage(image, input.Width, input.Height)
		byteImage, err = models.ConvertImageIntoByte(resizedImage, input.Extension)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	} else {
		byteImage, err = models.ConvertImageIntoByte(image, input.Extension)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}
	var filePath string
	if input.FilePath == "" {
		filePath = models.GenerateRandonFileName(input.Extension)
	} else {
		filePath = input.FilePath
	}

	publicURL, err := models.SaveAsS3(byteImage, input.Extension, user.AwsAccessKeyId, user.AwsSecretAccessKey, user.AwsRegion, user.AwsBucket, filePath)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	model := models.Image{
		PublicURL: publicURL,
	}

	context.JSON(http.StatusOK, model)
}
