package repository

import (
	"github.com/Zhiyenbek/sp-positions-main-service/config"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PositionRepository interface {
	GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error)
	GetPositionInterviews(publicID string, pageNum int, pageSize int) ([]models.InterviewResults, int, error)
	Exists(publicID string) (bool, error)
	GetPosition(publicID string) (*models.Position, error)
	CreatePosition(position *models.Position) (string, error)
	CreateSkillsForPosition(positionPublicID string, skills []string) error
	DeleteSkillsFromPosition(positionPublicID string, skills []string) error
	GetPositionsByCompany(companyID string, pageNum int, pageSize int, search string) ([]models.Position, int, error)
	GetPositionsByRecruiter(recruiterID string, pageNum int, pageSize int, search string) ([]models.Position, int, error)
	AddQuestionsToPosition(positionPublicID string, questions []*models.Question) ([]*models.Question, error)
}

type CompanyRepository interface {
	GetCompanyByRecruiterPublicID(recruiterPublicID string) (*models.Company, error)
}

type Repository struct {
	PositionRepository
	CompanyRepository
}

func New(db *pgxpool.Pool, cfg *config.Configs, log *zap.SugaredLogger) *Repository {
	return &Repository{
		PositionRepository: NewPositionRepository(db, cfg.DB, log),
		CompanyRepository:  NewCompanyRepository(db, cfg.DB, log),
	}
}
