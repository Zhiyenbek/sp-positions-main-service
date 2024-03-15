package repository

import (
	"github.com/Zhiyenbek/sp_positions_main_service/config"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PositionRepository interface {
	GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error)
}

type Repository struct {
	PositionRepository
}

func New(db *pgxpool.Pool, cfg *config.Configs, log *zap.SugaredLogger) *Repository {
	return &Repository{
		PositionRepository: NewPositionRepository(db, cfg.DB, log),
	}
}
