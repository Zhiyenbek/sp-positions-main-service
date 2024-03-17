package middleware

import (
	handler "github.com/Zhiyenbek/users-auth-service/internal/handler/http"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VerifyToken(tokenSecret string, log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := c.Cookie("access_token")
		if err != nil {
			log.Error("cookie not found")
			c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
			return
		}
		token, err := handler.ParseAuthToken(jwtToken, tokenSecret)
		if err != nil {
			log.Error("token is invalid", err)
			c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
			return
		}
		c.Set("role", token.Role)
		c.Set("public_id", token.PublicID)
		// Pass on to the next-in-chain
		c.Next()
	}
}
func sendResponse(status int, data interface{}, err error) gin.H {
	var errResponse gin.H
	if err != nil {
		errResponse = gin.H{
			"message": err.Error(),
		}
	} else {
		errResponse = nil
	}

	return gin.H{
		"data":   data,
		"status": status,
		"error":  errResponse,
	}
}
