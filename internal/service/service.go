package service

import (
	"github.com/Zhiyenbek/sp-positions-main-service/config"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/models"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/repository"
	"go.uber.org/zap"
)

type PositionService interface {
	GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error)
	GetPositionInterviews(publicID string, pageNum int, pageSize int) ([]models.Interview, int, error)
	Exists(publicID string) error
	GetPosition(publicID string) (*models.Position, error)
	CreatePosition(position *models.Position) (*models.Position, error)
	CreateSkillsForPosition(positionPublicID string, skills []string) error
	DeleteSkillsFromPosition(positionPublicID string, skills []string) error
	GetPositionsByCompany(companyID string, pageNum int, pageSize int, search string) ([]models.Position, int, error)
}
type Service struct {
	PositionService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		PositionService: NewPositionsService(repos, cfg, log),
	}
}
