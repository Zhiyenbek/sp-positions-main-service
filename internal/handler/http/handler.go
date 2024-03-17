package handler

import (
	"github.com/Zhiyenbek/sp-positions-main-service/config"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/service"
	"github.com/Zhiyenbek/users-auth-service/middleware"
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
	router.GET("/positions", h.GetPositions)
	router.GET("/positions/:position_public_id/interviews", h.GetPositionInterviews)
	router.GET("/position/:position_public_id", h.GetPosition)
	router.POST("/position", middleware.VerifyToken(h.cfg.Token.TokenSecret), h.CreatePosition)
	router.POST("/position/:position_public_id/skills", h.AddSkillsToPosition)
	router.DELETE("/position/:position_public_id/skills", h.DeleteSkillsFromPosition)
	// router.PUT("/position", middleware.VerifyToken(h.cfg.Token.TokenSecret), h.UpdatePosition)
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
