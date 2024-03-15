package repository

import (
	"context"

	"github.com/Zhiyenbek/sp_positions_main_service/config"
	"github.com/Zhiyenbek/sp_positions_main_service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type positionRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewPositionRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) PositionRepository {
	return &positionRepository{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
}

func (r *positionRepository) GetAllPositions(search string, pageNum, pageSize int) ([]models.Position, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	offset := (pageNum - 1) * pageSize

	// Query to retrieve positions for the current page
	query := `
		SELECT id, public_id, name, status
		FROM positions
		WHERE name ILIKE $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, "%"+search+"%", pageSize, offset)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving positions: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	positions := []models.Position{}
	for rows.Next() {
		position := models.Position{}
		err := rows.Scan(
			&position.ID,
			&position.PublicID,
			&position.Name,
			&position.Status,
		)
		if err != nil {
			r.logger.Errorf("Error occurred while scanning position rows: %v", err)
			return nil, 0, err
		}

		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Error occurred while iterating over position rows: %v", err)
		return nil, 0, err
	}

	// Query to retrieve the total count of positions
	countQuery := `
		SELECT COUNT(*) FROM positions WHERE name ILIKE $1
	`

	var totalCount int
	err = r.db.QueryRow(ctx, countQuery, "%"+search+"%").Scan(&totalCount)
	if err != nil {
		r.logger.Errorf("Error occurred while retrieving position count: %v", err)
		return nil, 0, err
	}

	return positions, totalCount, nil
}
