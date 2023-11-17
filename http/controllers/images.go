package controllers

import (
	"net/http"

	"github.com/OdairPianta/gryphon_api/http/requests"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/gin-gonic/gin"
)

// @Summary Get image
// @Schemes
// @Description Get image
// @Tags images
// @Accept json
// @Produce json
// @Param   id     path    int     true        "Image ID"
// @Success 200 {object} models.Image "ok"
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router /images/{id} [get]
// @Security Bearer
func FindImage(context *gin.Context) {
	var model models.Image
	err := models.DB.Where("id = ?", context.Param("id")).First(&model).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Image not found!"})
		return
	}

	context.JSON(http.StatusOK, model)
}

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

	publicURL, err := models.SaveAsS3(input.Base64, input.Extension, user.AwsAccessKeyId, user.AwsSecretAccessKey, user.AwsRegion)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	model := models.Image{
		PublicURL: publicURL,
	}

	context.JSON(http.StatusOK, model)
}

// @Summary		Delete an image
// @Description	Delete by image ID
// @Tags			images
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"ID"	Format(int64)
// @Success		200	{object}	models.Image
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/images/{id} [delete]
// @Security Bearer
func DeleteImage(context *gin.Context) {
	var model models.Image
	err := models.DB.Where("id = ?", context.Param("id")).First(&model).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Image not found!"})
		return
	}

	result := models.DB.Delete(&model)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, model)
}

// @Summary		Delete an image by public url
// @Description	Delete by image public url
// @Tags			images
// @Accept			json
// @Produce		json
// @Param			public_url	path		string	true	"Public URL"
// @Success		200	{object}	string
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/images/public_url/{public_url} [delete]
// @Security Bearer
func DeleteImageByPublicURL(context *gin.Context) {
	var model models.Image
	err := models.DB.Where("public_url = ?", context.Param("public_url")).First(&model).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Image not found!"})
		return
	}

	result := models.DB.Delete(&model)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, model)
}
