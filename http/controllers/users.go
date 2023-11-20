package controllers

import (
	"net/http"
	"time"

	"github.com/OdairPianta/gryphon_api/http/requests"
	"github.com/OdairPianta/gryphon_api/models"
	"github.com/OdairPianta/gryphon_api/notifications"
	"github.com/OdairPianta/gryphon_api/services/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Get users
// @Schemes
// @Description Get users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} []models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router /users [get]
// @Security Bearer
func FindUsers(context *gin.Context) {
	var users []models.User
	models.DB.Find(&users)

	context.JSON(http.StatusOK, users)
}

// @Summary Get user
// @Schemes
// @Description Get user
// @Tags users
// @Accept json
// @Produce json
// @Param   id     path    int     true        "User ID"
// @Success 200 {object} models.User "ok"
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router /users/{id} [get]
// @Security Bearer
func FindUser(context *gin.Context) {
	var user models.User
	err := models.DB.Where("id = ?", context.Param("id")).First(&user).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado!"})
		return
	}

	context.JSON(http.StatusOK, user)
}

// @Summary		Add an user
// @Description	add by json user
// @Tags			users
// @Accept		json
// @Produce		json
// @Param			user	body		requests.CreateUserInput	true	"Add user"
// @Success		200		{object}	models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/users [post]
// @Security Bearer
func CreateUser(context *gin.Context) {
	// Validate input
	var input requests.CreateUserInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, errEncryptPass := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if errEncryptPass != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao criptografar a senha!"})
		return
	}

	// Create user
	user := models.User{
		Name:            input.Name,
		Email:           input.Email,
		EmailVerifiedAt: time.Now(),
		Password:        string(hashedPassword),
		FCMToken:        input.FCMToken,
	}
	resultInsert := models.DB.Create(&user)

	if resultInsert.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": resultInsert.Error.Error()})
		return
	}

	token, errGenerateToken := token.GenerateToken(user.ID)
	if errGenerateToken != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errGenerateToken.Error()})
		return
	}
	user.Token = token

	context.JSON(http.StatusOK, user)
}

// @Summary		Update an user
// @Description	Update by json user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id		path		int					true	"ID"
// @Param			user	body		requests.UpdateUserInput	true	"Update user"
// @Success		200		{object}	models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/users/{id} [put]
// @Security Bearer
func UpdateUser(context *gin.Context) {
	// Get model if exist
	var user models.User
	err := models.DB.Where("id = ?", context.Param("id")).First(&user).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado!"})
		return
	}

	// Validate input
	var input requests.UpdateUserInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := models.DB.Model(&user).Updates(input)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

// @Summary		Update an user fcm token
// @Description	Update by json user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id		path		int					true	"ID"
// @Param			user	body		requests.UpdateUserFcmTokenInput	true	"Update user fcm token"
// @Success		200		{object}	models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/users/{id}/update_fcm_token [put]
// @Security Bearer
func UpdateFcmToken(context *gin.Context) {
	// Get model if exist
	var user models.User
	err := models.DB.Where("id = ?", context.Param("id")).First(&user).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado!"})
		return
	}

	// Validate input
	var input requests.UpdateUserFcmTokenInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result := models.DB.Model(&user).Updates(input)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

// @Summary		Delete an user
// @Description	Delete by user ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"ID"	Format(int64)
// @Success		200	{object}	models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/users/{id} [delete]
// @Security Bearer
func DeleteUser(context *gin.Context) {
	// Get model if exist
	var user models.User
	err := models.DB.Where("id = ?", context.Param("id")).First(&user).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado!"})
		return
	}

	result := models.DB.Delete(&user)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

// @Summary		Login an user
// @Description	Login by json user
// @Tags			users
// @Accept		json
// @Produce		json
// @Param			user	body		requests.LoginInput	true	"Add user"
// @Success		200		{object}	models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/login [post]
func Login(context *gin.Context) {

	var input requests.LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := models.User{}

	user.Email = input.Email
	user.Password = input.Password

	token, loginErr := models.LoginCheck(user.Email, user.Password)

	if loginErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email ou senha incorretos" + loginErr.Error()})
		return
	}
	//get user by email
	errGetUser := models.DB.Where("email = ?", user.Email).First(&user).Error
	if errGetUser != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado" + errGetUser.Error()})
		return
	}
	user.Token = token

	context.JSON(http.StatusOK, user)
}

// @Summary		Forgot password
// @Description	Forgot password by json user
// @Tags			users
// @Accept		json
// @Produce		json
// @Param			user	body		requests.ForgotPasswordInput	true	"Forgot password"
// @Success		200		{object}	string
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/forgot_password [post]
func ForgotPassword(context *gin.Context) {
	var input requests.ForgotPasswordInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := models.User{}

	user.Email = input.Email

	err := models.DB.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado" + err.Error()})
		return
	}

	token, errGenerateToken := token.GenerateToken(user.ID)
	if errGenerateToken != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errGenerateToken.Error()})
		return
	}
	user.ResetPasswordToken = token

	result := models.DB.Model(&user).Updates(user)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	forgotPasswordError := notifications.ForgotPassword(user)
	if forgotPasswordError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": forgotPasswordError.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Email enviado com sucesso!"})
}

// @Summary		Recover password
// @Description	Recover password by json user
// @Tags			users
// @Accept		json
// @Produce		json
// @Param			user	body		requests.RecoverPasswordInput	true	"Recover password"
// @Success		200		{object}	models.User
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Router			/recover_password [post]
func RecoverPassword(context *gin.Context) {
	var input requests.RecoverPasswordInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := models.User{}

	user.ResetPasswordToken = input.ResetPasswordToken
	user.Password = input.Password

	err := models.DB.Where("reset_password_token = ?", user.ResetPasswordToken).First(&user).Error
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Token não encontrado" + err.Error()})
		return
	}

	hashedPassword, errEncryptPass := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if errEncryptPass != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao criptografar senha!"})
		return
	}

	user.Password = string(hashedPassword)
	user.ResetPasswordToken = ""

	result := models.DB.Model(&user).Updates(map[string]interface{}{"password": user.Password, "reset_password_token": user.ResetPasswordToken})
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	token, loginErr := models.LoginCheck(user.Email, input.Password)

	if loginErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "O email ou a senha estão incorretos." + loginErr.Error()})
		return
	}

	user.Token = token

	context.JSON(http.StatusOK, user)
}
