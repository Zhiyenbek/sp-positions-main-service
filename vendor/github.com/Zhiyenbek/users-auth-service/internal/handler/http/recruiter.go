package handler

import (
	"errors"
	"net/http"

	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *handler) RecruiterSignUp(c *gin.Context) {
	req := &models.RecruiterSignUpRequest{}

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

	err := h.service.CreateRecruiter(req)
	if err != nil {
		var errMsg error
		var code int
		switch {
		case errors.Is(err, models.ErrCompanyDoesntExists):
			errMsg = models.ErrCompanyDoesntExists
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

func (h *handler) RecruiterSignIn(c *gin.Context) {
	req := &models.UserSignInRequest{}
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		h.logger.Errorf("ERROR: invalid input, some fields are incorrect: %s\n", err.Error())
		c.AbortWithStatusJSON(400, sendResponse(-1, nil, models.ErrInvalidInput))
		return
	}

	if !validatePassword(req.Password) {
		h.logger.Error("failed to validate password")
		c.JSON(400, sendResponse(-1, nil, models.ErrWrongCredential))
		return
	}
	var errMsg error
	tokens, err := h.service.AuthService.RecruiterLogin(req)
	if err != nil {
		h.logger.Errorf("Error occurred while login: %v", err)
		switch {
		case errors.Is(err, models.ErrWrongCredential):
			errMsg = models.ErrWrongCredential
		default:
			c.JSON(500, sendResponse(-1, nil, models.ErrInternalServer))
			return
		}
		c.JSON(200, sendResponse(-1, nil, errMsg))
		return
	}
	c.SetCookie("access_token", tokens.AccessToken.TokenValue, int(tokens.AccessToken.TTL.Seconds()), "/", h.cfg.Token.Access.Domain, true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken.TokenValue, int(tokens.RefreshToken.TTL.Seconds()), "/refresh-token", h.cfg.Token.Refresh.Domain, true, true)
	c.JSON(http.StatusOK, sendResponse(0, nil, nil))
}
