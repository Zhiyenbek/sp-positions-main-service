package handler

import (
	"unicode"

	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *handler) TestAuth(c *gin.Context) {
	c.JSON(200, sendResponse(0, "YOU ARE AUTHORIZED", nil))
}

func (h *handler) RefreshToken(c *gin.Context) {
	rtToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
		return
	}
	tokens, err := h.service.AuthService.RefreshToken(rtToken)
	if err != nil {
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}
	c.SetCookie("access_token", tokens.AccessToken.TokenValue, int(tokens.AccessToken.TTL.Seconds()), "/", h.cfg.Token.Access.Domain, true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken.TokenValue, int(tokens.RefreshToken.TTL.Seconds()), "/refresh-token", h.cfg.Token.Refresh.Domain, true, true)

	c.JSON(200, sendResponse(0, nil, nil))
}

// validatePassword - function that validates password. Password being validated by these requirements:
// 1.Password must have upper case characters
// 2.Password must have special characters
// function returns boolean
func validatePassword(pass string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
