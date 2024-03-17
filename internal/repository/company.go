package repository

import (
	"context"

	"github.com/Zhiyenbek/sp_positions_main_service/config"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type companyRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewCompanyRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) CompanyRepository {
	return &companyRepository{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
}

func (r *companyRepository) GetCompanyByRecruiterPublicID(recruiterPublicID string) (*models.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	query := `SELECT c.name, c.public_id, c.logo, c.description 
		FROM companies AS c
		INNER JOIN recruiters r ON r.company_public_id = c.public_id
		WHERE r.public_id = $1
		GROUP BY c.name, c.public_id, c.logo, c.description `

	row := r.db.QueryRow(ctx, query, recruiterPublicID)

	company := &models.Company{}
	err := row.Scan(&company.Name, &company.PublicID, &company.Logo, &company.Description)
	if err != nil {
		r.logger.Errorf("Error occurred while fetching company: %v", err)
		return nil, err
	}

	return company, nil
}
