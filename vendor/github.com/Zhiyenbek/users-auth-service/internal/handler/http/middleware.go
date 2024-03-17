package handler

import (
	"fmt"
	"time"

	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ParseAuthToken(tokenString string, tokenSecret string) (*models.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v %w", err, models.ErrInvalidInput)
	}
	if claims, ok := token.Claims.(*models.JwtUserClaims); ok && token.Valid {
		token := &models.Token{
			PublicID:   claims.PublicID,
			TokenValue: tokenString,
			TTL:        time.Duration(claims.ExpiresAt),
		}
		return token, nil
	}
	return nil, fmt.Errorf("could not parse token: %w", models.ErrInvalidToken)
}

func (h *handler) VerifyToken(c *gin.Context) {
	jwtToken, err := c.Cookie("access_token")
	if err != nil {
		h.logger.Error(err)
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
		return
	}
	_, err = ParseAuthToken(jwtToken, h.cfg.Token.Access.TokenSecret)
	if err != nil {
		h.logger.Error(err)
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
		return
	}
	// Pass on to the next-in-chain
	c.Next()
}
