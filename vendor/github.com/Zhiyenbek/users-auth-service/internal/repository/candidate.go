package repository

import (
	"context"

	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// candidateRepository represents the repository for managing candidates in the database.
type candidateRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

// NewCandidateRepository creates a new instance of candidateRepository.
func NewCandidateRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) CandidateRepository {
	return &candidateRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *candidateRepository) CreateCandidate(candidate *models.CandidateSignUpRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	var user_id, candidate_id int64
	var user_public_id uuid.UUID
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Errorf("Error occurred while creating recruiter in users: %v", err)
		return err
	}

	query := `INSERT INTO users 
				(first_name, last_name)
			VALUES
				($1, $2)
			RETURNING id, public_id`

	err = tx.QueryRow(ctx, query, candidate.FirstName, candidate.LastName).Scan(&user_id, &user_public_id)
	if err != nil {
		r.logger.Errorf("Error occurred while creating candidate in users: %v", err)

		errTX := tx.Rollback(ctx)
		if errTX != nil {
			r.logger.Errorf("ERROR: transaction: %s", errTX)
		}
		return err
	}

	query = `INSERT INTO candidates (public_id, current_position, resume, bio, education) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err = tx.QueryRow(ctx, query, user_public_id, candidate.CurrentPosition, candidate.Resume, candidate.Bio, candidate.Education).Scan(&candidate_id)
	if err != nil {
		r.logger.Errorf("Error occurred while creating candidate: %v", err)

		errTX := tx.Rollback(ctx)
		if errTX != nil {
			r.logger.Errorf("ERROR: transaction: %s", errTX)
		}
		return err
	}

	if candidate.Skills != nil && len(candidate.Skills) > 0 {
		// Assuming candidate.Skills is a slice of skill names
		for _, skill := range candidate.Skills {
			var skillID int

			// Check if the skill already exists
			query := `SELECT id FROM skills WHERE name = $1;`
			err := tx.QueryRow(ctx, query, skill).Scan(&skillID)
			if err != nil {
				// Skill doesn't exist, insert it and get the ID
				query := `INSERT INTO skills (name) VALUES ($1) RETURNING id;`
				err = tx.QueryRow(ctx, query, skill).Scan(&skillID)
				if err != nil {
					r.logger.Errorf("Error occurred while creating skill for candidate: %v", err)

					errTX := tx.Rollback(ctx)
					if errTX != nil {
						r.logger.Errorf("ERROR: transaction: %s", errTX)
					}
					return err
				}
			}

			// Associate the skill with the candidate
			query = `INSERT INTO candidate_skills (candidate_id, skill_id) VALUES ($1, $2);`
			_, err = tx.Exec(ctx, query, candidate_id, skillID)
			if err != nil {
				r.logger.Errorf("Error occurred while associating skill with candidate: %v", err)

				errTX := tx.Rollback(ctx)
				if errTX != nil {
					r.logger.Errorf("ERROR: transaction: %s", errTX)
				}
				return err
			}
		}
	}

	query = `INSERT INTO auth (user_id, login, password) VALUES ($1, $2, $3);`

	_, err = tx.Exec(ctx, query, user_id, candidate.Login, candidate.Password)
	if err != nil {
		r.logger.Errorf("Error occurred while creating authentication info: %v", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Errorf("Error occurred while committing transaction: %v", err)

		errTX := tx.Rollback(ctx)
		if errTX != nil {
			r.logger.Errorf("ERROR: transaction error: %s", errTX)
		}
		return err
	}

	return nil
}

func (r *candidateRepository) Exists(publicID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM candidates WHERE public_id = $1)`

	err := r.db.QueryRow(ctx, query, publicID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return false, err
	}

	return exists, nil
}
