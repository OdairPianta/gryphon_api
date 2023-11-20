package controllers

import (
	"net/http"

	"github.com/OdairPianta/gryphon_api/http/requests"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/gin-gonic/gin"
)

// @Summary		Add base64 image
// @Description	Add base64 image
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

	publicURL, err := models.SaveAsS3(input.Base64, input.Extension, user.AwsAccessKeyId, user.AwsSecretAccessKey, user.AwsRegion, user.AwsBucket)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	model := models.Image{
		PublicURL: publicURL,
	}

	context.JSON(http.StatusOK, model)
}
