package controllers

import (
	"net/http"

	"github.com/OdairPianta/gryphon_api/models"
	"github.com/gin-gonic/gin"
)

// @BasePath /api

// @Summary Get api status
// @Schemes
// @Description Get api status
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} models.APIStatus "ok"
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router /status [get]
func FindStatus(context *gin.Context) {
	var clientStatus string = "pending status"
	var apiStatus = models.APIStatus{
		ClientStatus: clientStatus,
	}

	context.JSON(http.StatusOK, apiStatus)
}
