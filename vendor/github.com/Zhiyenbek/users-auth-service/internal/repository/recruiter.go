package repository

import (
	"context"
	"errors"

	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type recruiterRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewRecruiterRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) RecruiterRepository {
	return &recruiterRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *recruiterRepository) CreateRecruiter(recruiter *models.RecruiterSignUpRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()
	var user_id int64
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

	err = tx.QueryRow(ctx, query, recruiter.FirstName, recruiter.LastName).Scan(&user_id, &user_public_id)
	if err != nil {
		r.logger.Errorf("Error occurred while creating recruiter in users: %v", err)

		errTX := tx.Rollback(ctx)
		if errTX != nil {
			r.logger.Errorf("ERROR: transaction: %s", errTX)
		}
		return err
	}
	query = `SELECT public_id FROM companies WHERE name = $1;`
	err = tx.QueryRow(ctx, query, recruiter.CompanyPublicID).Scan(&recruiter.CompanyPublicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			query := `INSERT INTO companies (name)
			VALUES
				($1)
			RETURNING public_id;`

			err = tx.QueryRow(ctx, query, recruiter.CompanyName).Scan(&recruiter.CompanyPublicID)
			if err != nil {
				r.logger.Errorf("Error occurred while creating recruiters in companies: %v", err)

				errTX := tx.Rollback(ctx)
				if errTX != nil {
					r.logger.Errorf("ERROR: transaction: %s", errTX)
				}
				return err
			}
		} else {
			r.logger.Errorf("Error occurred while checking company existence: %v", err)
			errTX := tx.Rollback(ctx)
			if errTX != nil {
				r.logger.Errorf("ERROR: transaction: %s", errTX)
			}
			return err
		}
	}

	query = `INSERT INTO auth (user_id, login, password) VALUES ($1, $2, $3);`

	_, err = tx.Exec(ctx, query, user_id, recruiter.Login, recruiter.Password)
	if err != nil {
		r.logger.Errorf("Error occurred while creating authentication info: %v %d %s %s", err, user_id, recruiter.Login, recruiter.Password)
		return err
	}

	query = `INSERT INTO recruiters
				(public_id, company_public_id)
			 VALUES ($1, $2)`

	_, err = tx.Exec(ctx, query, user_public_id, recruiter.CompanyPublicID)

	if err != nil {
		r.logger.Errorf("Error occurred while creating recruiters: %v", err)
		errTX := tx.Rollback(ctx)
		if errTX != nil {
			r.logger.Errorf("ERROR: transaction: %s", errTX)
		}
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

func (r *recruiterRepository) Exists(publicID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM recruiters WHERE public_id = $1)`

	err := r.db.QueryRow(ctx, query, publicID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return false, err
	}

	return exists, nil
}
