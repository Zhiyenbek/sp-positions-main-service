package handler

import (
	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handler struct {
	service *service.Service
	cfg     *config.Configs
	logger  *zap.SugaredLogger
}

type Handler interface {
	InitRoutes() *gin.Engine
}

func New(services *service.Service, logger *zap.SugaredLogger, cfg *config.Configs) Handler {
	return &handler{
		service: services,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/recruiter/sign-in", h.RecruiterSignIn)
	router.POST("/candidate/sign-in", h.CandidateSignIn)
	router.POST("/recruiter/sign-up", h.RecruiterSignUp)
	router.POST("/candidate/sign-up", h.CandidateSignUp)
	router.POST("/refresh-token", h.RefreshToken)

	router.POST("/verify", h.VerifyToken, h.TestAuth)

	return router
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
