package routes

import (
	"os"

	docs "github.com/OdairPianta/gryphon_api/docs"
	"github.com/OdairPianta/gryphon_api/http/controllers"
	"github.com/OdairPianta/gryphon_api/http/middlewares"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequest() {
	r := gin.Default()
	InitRoutes(r)
	r.Run(":" + os.Getenv("APP_PORT"))
}

func InitRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middlewares.JSONLogMiddleware())

	public := r.Group("/api")
	docs.SwaggerInfo.BasePath = "/api"
	public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	public.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	public.GET("/status", controllers.FindStatus)

	public.POST("/login", controllers.Login)
	public.POST("/forgot_password", controllers.ForgotPassword)
	public.POST("/recover_password", controllers.RecoverPassword)

	protected := r.Group("/api")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.POST("/images/base64/create", controllers.CreateBase64Image)

	protected.GET("users", controllers.FindUsers)
	protected.POST("users", controllers.CreateUser)
	protected.GET("/users/:id", controllers.FindUser)
	protected.PUT("/users/:id", controllers.UpdateUser)
	protected.PUT("/users/:id/update_fcm_token", controllers.UpdateFcmToken)
	protected.DELETE("/users/:id", controllers.DeleteUser)

	return r
}
