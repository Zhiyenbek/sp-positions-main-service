package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type authRepository struct {
	db     *pgxpool.Pool
	cfg    *config.DBConf
	logger *zap.SugaredLogger
}

func NewAuthRepository(db *pgxpool.Pool, cfg *config.DBConf, logger *zap.SugaredLogger) AuthRepository {
	return &authRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *authRepository) GetUserInfoByLogin(login string) (string, string, error) {
	var password string
	var ID uuid.UUID
	timeout := r.cfg.TimeOut
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	query := `SELECT u.public_id, a.password
	FROM users as u
	JOIN auth as a ON u.id = a.user_id
	WHERE a.login= $1`
	if err := r.db.QueryRow(ctx, query, login).Scan(&ID, &password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", fmt.Errorf("%w: error occurred while getting password_hash from db: auth does not exist", models.ErrWrongCredential)
		}
		return "", "", fmt.Errorf("%w: error occurred while getting password_hash from db: %v", models.ErrInternalServer, err)
	}
	return password, ID.String(), nil

}

func (r *authRepository) Exists(login string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.TimeOut)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth WHERE login = $1)`
	fmt.Println(r.logger)
	err := r.db.QueryRow(ctx, query, login).Scan(&exists)
	if err != nil {
		r.logger.Errorf("Error occurred while checking user existence: %v", err)
		return false, err
	}

	return exists, nil
}
