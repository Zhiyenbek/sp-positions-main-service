package service

import (
	"encoding/json"

	"github.com/Zhiyenbek/sp-positions-main-service/config"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/models"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/repository"
	"go.uber.org/zap"
)

type positionsService struct {
	cfg          *config.Configs
	logger       *zap.SugaredLogger
	positionRepo repository.PositionRepository
	companyRepo  repository.CompanyRepository
}

func NewPositionsService(repo *repository.Repository, cfg *config.Configs, logger *zap.SugaredLogger) PositionService {
	return &positionsService{
		positionRepo: repo.PositionRepository,
		companyRepo:  repo.CompanyRepository,
		cfg:          cfg,
		logger:       logger,
	}
}

func (p *positionsService) GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error) {
	return p.positionRepo.GetAllPositions(search, pageNum, pageSize)
}

func (p *positionsService) GetPosition(publicID string) (*models.Position, error) {
	return p.positionRepo.GetPosition(publicID)
}

func (p *positionsService) GetPositionInterviews(publicID string, pageNum int, pageSize int) ([]*models.Interview, int, error) {
	interviewRawResult, count, err := p.positionRepo.GetPositionInterviews(publicID, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	res := make([]*models.Interview, 0)
	for _, r := range interviewRawResult {
		if r.Result != nil {
			result := &models.Result{}
			interview := &models.Interview{}
			err = json.Unmarshal(r.Result, &result)
			if err != nil {
				p.logger.Error(err)
				return nil, 0, err
			}
			interview.PublicID = r.PublicID
			interview.Result = *result
			res = append(res, interview)
		} else {
			res = append(res, &models.Interview{
				PublicID: r.PublicID,
			})
		}
	}
	return res, count, nil
}

func (p *positionsService) CreatePosition(position *models.Position) (*models.Position, error) {
	publicID, err := p.positionRepo.CreatePosition(position)
	if err != nil {
		return nil, err
	}
	company, err := p.companyRepo.GetCompanyByRecruiterPublicID(*position.RecruiterPublicID)
	if err != nil {
		return nil, err
	}
	position.PublicID = &publicID
	position.Company = company
	return position, nil
}

func (p *positionsService) Exists(publicID string) error {
	exists, err := p.positionRepo.Exists(publicID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrPositionNotFound
	}
	return nil
}

func (p *positionsService) CreateSkillsForPosition(positionPublicID string, skills []string) error {
	return p.positionRepo.CreateSkillsForPosition(positionPublicID, skills)
}

func (p *positionsService) DeleteSkillsFromPosition(positionPublicID string, skills []string) error {
	return p.positionRepo.DeleteSkillsFromPosition(positionPublicID, skills)
}

func (p *positionsService) GetPositionsByCompany(companyID string, pageNum int, pageSize int, search string) ([]models.Position, int, error) {
	return p.positionRepo.GetPositionsByCompany(companyID, pageNum, pageSize, search)
}

func (p *positionsService) GetPositionsByRecruiter(recruiterID string, pageNum int, pageSize int, search string) ([]models.Position, int, error) {
	return p.positionRepo.GetPositionsByRecruiter(recruiterID, pageNum, pageSize, search)
}

func (p *positionsService) AddQuestionsToPosition(positionPublicID string, questions []*models.Question) ([]*models.Question, error) {
	return p.positionRepo.AddQuestionsToPosition(positionPublicID, questions)
}

func (p *positionsService) GetPositionQuestions(positionPublicID string) ([]*models.Question, error) {
	return p.positionRepo.GetPositionQuestions(positionPublicID)
}

func (p *positionsService) CreateInterview(positionPublicID, candidatePublicID string) (string, error) {
	return p.positionRepo.CreateInterview(positionPublicID, candidatePublicID)
}
