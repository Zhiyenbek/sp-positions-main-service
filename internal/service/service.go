package service

import (
	"github.com/Zhiyenbek/sp_positions_main_service/config"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/models"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/repository"
	"go.uber.org/zap"
)

type PositionService interface {
	GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error)
}
type Service struct {
	PositionService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		PositionService: NewPositionsService(repos, cfg, log),
	}
}
