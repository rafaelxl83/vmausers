package middlewares

import (
	"net/http"
	"vmausers/auth"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Authorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			log.Error("Authorization: Bad Request [request does not contain an access token]")
			context.JSON(http.StatusBadRequest, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Errorf("Authorization: Unauthorized [v%]", err)
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		log.Info("Authorized")
		context.Next()
	}
}
