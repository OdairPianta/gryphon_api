package middlewares

import (
	"net/http"

	"github.com/OdairPianta/gryphon_api/services/token"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := token.TokenValid(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "Você não está logado"})
			context.Abort()
			return
		}
		context.Next()
	}
}
