package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/OdairPianta/gryphon_api/models"
	"github.com/OdairPianta/gryphon_api/routes"
	"github.com/OdairPianta/gryphon_api/services/token"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func setupDatabase() (*gorm.DB, func()) {
	models.ConnectDatabase()
	tx := models.DB.Begin()

	cleanup := func() {
		tx.Rollback()
	}

	return tx, cleanup
}

func routesSetup() *gin.Engine {
	r := gin.Default()
	routes.InitRoutes(r)
	return r
}

func initUser(withFieldsToIgnore ...string) (user models.User, stringToken string) {
	user = models.User{}
	withFieldsToIgnore = append(withFieldsToIgnore, "DeletedAt")
	err := faker.FakeData(&user, options.WithFieldsToIgnore(withFieldsToIgnore...))
	if err != nil {
		fmt.Println(err)
		return
	}
	user.ID = 0
	result := models.DB.Create(&user)
	if result.Error != nil {
		fmt.Println(err)
		return
	}

	token, errToken := token.GenerateToken(user.ID)
	if errToken != nil {
		fmt.Println(err)
		return
	}
	return user, token
}

func TestUserCreateExistInDatabaseAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	_, token := initUser()

	model := models.User{}

	err := faker.FakeData(&model, options.WithFieldsToIgnore("DeletedAt", "ID", "RememberToken"))
	assert.Nil(t, err)

	model.Email = faker.Email()
	jsonModel, err := json.Marshal(model)
	assert.Nil(t, err)

	request, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonModel))
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var result map[string]string
	_ = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.NotNil(t, result, "response body must be a valid json")

	assert.NotEmpty(t, result["token"], "token must be not empty")
}

func TestUserLoginIsPerformedAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	user, _ := initUser()

	body := []byte(`{"email": "` + user.Email + `", "password": "123456"}`)

	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var result map[string]string
	_ = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.NotNil(t, result, "response body must be a valid json")

	assert.NotEmpty(t, result["token"], "token must be not empty")
}

func TestFindUsersReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	_, token := initUser()

	request, _ := http.NewRequest("GET", "/api/users", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var result []models.User
	_ = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.NotNil(t, result, "response body must be a valid json")
}

func TestUserFindExistInDatabaseAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	_, token := initUser()

	model := models.User{}
	err := faker.FakeData(&model, options.WithFieldsToIgnore("DeletedAt", "ID"))
	assert.Nil(t, err)

	result := models.DB.Create(&model)
	assert.Nil(t, result.Error)

	request, _ := http.NewRequest("GET", "/api/users/"+fmt.Sprint(model.ID), nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var resultModel models.User
	_ = json.Unmarshal(recorder.Body.Bytes(), &resultModel)
	assert.NotNil(t, resultModel, "response body must be a valid json")
}

func TestUserUpdateExistInDatabaseAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	_, token := initUser()

	model := models.User{}
	err := faker.FakeData(&model, options.WithFieldsToIgnore("DeletedAt", "ID"))
	if err != nil {
		fmt.Println(err)
	}

	result := models.DB.Create(&model)
	if result.Error != nil {
		fmt.Println(err)
		return
	}

	model.Name = faker.Name()
	jsonModel, err := json.Marshal(model)
	if err != nil {
		fmt.Println(err)
		return
	}

	request, _ := http.NewRequest("PUT", "/api/users/"+fmt.Sprint(model.ID), bytes.NewBuffer(jsonModel))
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var resultModel models.User
	_ = json.Unmarshal(recorder.Body.Bytes(), &resultModel)
	assert.NotNil(t, resultModel, "response body must be a valid json")
}

func TestUserUpdateFcmTokenExistInDatabaseAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	_, token := initUser()

	model := models.User{}
	err := faker.FakeData(&model, options.WithFieldsToIgnore("DeletedAt", "ID"))
	if err != nil {
		fmt.Println(err)
	}

	result := models.DB.Create(&model)
	if result.Error != nil {
		fmt.Println(err)
		return
	}

	body := []byte(`{"fcm_token": "new_fcm_token"}`)

	request, _ := http.NewRequest("PUT", "/api/users/"+fmt.Sprint(model.ID)+"/update_fcm_token", bytes.NewBuffer(body))
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var resultModel models.User
	_ = json.Unmarshal(recorder.Body.Bytes(), &resultModel)
	assert.NotNil(t, resultModel, "response body must be a valid json")
	//assert fcm token is equal new_fcm_token
	assert.Equal(t, "new_fcm_token", resultModel.FCMToken, "fcm token must be equal new_fcm_token")
}

func TestUserDeleteIsPerformedAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	_, token := initUser()

	model := models.User{}
	err := faker.FakeData(&model, options.WithFieldsToIgnore("DeletedAt", "ID"))
	if err != nil {
		fmt.Println(err)
	}

	result := models.DB.Create(&model)
	if result.Error != nil {
		fmt.Println(err)
		return
	}

	jsonModel, err := json.Marshal(model)
	if err != nil {
		fmt.Println(err)
		return
	}

	request, _ := http.NewRequest("DELETE", "/api/users/"+fmt.Sprint(model.ID), bytes.NewBuffer(jsonModel))
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var resultModel models.User
	_ = json.Unmarshal(recorder.Body.Bytes(), &resultModel)
	assert.NotNil(t, resultModel, "response body must be a valid json")
}

func TestUserRecoverPasswordIsPerformedAndReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	user, _ := initUser("Token")

	body := []byte(`{"reset_password_token": "` + user.ResetPasswordToken + `", "password": "123456"}`)

	req, _ := http.NewRequest("POST", "/api/recover_password", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")
}

func TestStatusFindReturnCorrectData(t *testing.T) {
	setupDatabase()
	router := routesSetup()

	req, _ := http.NewRequest("GET", "/api/status", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")

	var result map[string]string
	_ = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.NotNil(t, result, "response body must be a valid json")
}

func TestUserCreateBase64Image(t *testing.T) {
	setupDatabase()
	router := routesSetup()
	var user models.User
	models.DB.First(&user, "email = ?", os.Getenv("USER_EMAIL"))
	token, errToken := token.GenerateToken(user.ID)
	if errToken != nil {
		fmt.Println(errToken)
		return
	}
	data := []byte(`{
		"base64": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdASFxOQERXRTc4UG1RV19iZ2hnPk1xeXBkeFxlZ2P/2wBDARESEhgVGC8aGi9jQjhCY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2P/wAARCAAFAAYDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwDFoooryz7w/9k=",
		"extension": "jpg"
	}`)
	request, _ := http.NewRequest("POST", "/api/images/base64/create", bytes.NewBuffer(data))
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	fmt.Println(recorder.Body.String())

	assert.Equal(t, http.StatusOK, recorder.Code, "OK response is expected")
}
