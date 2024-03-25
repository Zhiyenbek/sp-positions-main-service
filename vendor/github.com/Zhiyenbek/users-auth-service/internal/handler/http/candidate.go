package handler

import (
	"errors"
	"net/http"

	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *handler) CandidateSignUp(c *gin.Context) {
	req := &models.CandidateSignUpRequest{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("failed to parse request body when signing up candidate. %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	if !validatePassword(req.Password) {
		h.logger.Error("failed to validate password")
		c.JSON(400, sendResponse(-1, nil, models.ErrWrongCredential))
		return
	}
	err := h.service.CreateCandidate(req)
	if err != nil {
		var errMsg error
		var code int
		switch {
		case errors.Is(err, models.ErrUsernameExists):
			errMsg = models.ErrUsernameExists
			code = http.StatusBadRequest
		default:
			errMsg = models.ErrInternalServer
			code = http.StatusInternalServerError
		}
		c.JSON(code, sendResponse(-1, nil, errMsg))
		return
	}

	c.JSON(http.StatusCreated, sendResponse(0, nil, nil))
}

func (h *handler) CandidateSignIn(c *gin.Context) {
	req := &models.UserSignInRequest{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("ERROR: invalid input, some fields are incorrect: %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	switch {
	case !validatePassword(req.Password):
		h.logger.Error("invalid password")
		c.JSON(400, sendResponse(-1, nil, models.ErrInvalidPasswordFormat))
		return
	}

	tokens, err := h.service.AuthService.CandidateLogin(req)
	if err != nil {
		h.logger.Errorf("Error occurred while login: %v", err)
		switch {
		case errors.Is(err, models.ErrWrongCredential):
			c.JSON(http.StatusBadRequest, sendResponse(-1, nil, models.ErrWrongCredential))
		default:
			c.JSON(http.StatusInternalServerError, sendResponse(-1, nil, models.ErrInternalServer))
			return
		}
	}
	c.SetCookie("access_token", tokens.AccessToken.TokenValue, int(tokens.AccessToken.TTL.Seconds()), "/", h.cfg.Token.Access.Domain, true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken.TokenValue, int(tokens.RefreshToken.TTL.Seconds()), "/refresh-token", h.cfg.Token.Refresh.Domain, true, true)
	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}

func (h *handler) SignOut(c *gin.Context) {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		h.logger.Error(err)
		c.AbortWithStatusJSON(401, sendResponse(-1, nil, models.ErrInvalidToken))
		return
	}
	// Clear the access_token and refresh_token cookies
	c.SetCookie("access_token", "", -1, "/", h.cfg.Token.Access.Domain, true, true)
	c.SetCookie("refresh_token", "", -1, "/refresh-token", h.cfg.Token.Refresh.Domain, true, true)

	err = h.service.AuthService.SignOut(cookie)
	if err != nil {
		c.AbortWithStatusJSON(500, sendResponse(-1, nil, models.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}
