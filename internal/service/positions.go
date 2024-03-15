package service

import (
	"github.com/Zhiyenbek/sp_positions_main_service/config"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/models"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/repository"
	"go.uber.org/zap"
)

type positionsService struct {
	cfg          *config.Configs
	logger       *zap.SugaredLogger
	positionRepo repository.PositionRepository
}

func NewPositionsService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) PositionService {
	return &positionsService{
		positionRepo: repo.PositionRepository,
		cfg:          cfg,
		logger:       logger,
	}
}

func (p positionsService) GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error) {
	return p.positionRepo.GetAllPositions(search, pageNum, pageSize)
}
