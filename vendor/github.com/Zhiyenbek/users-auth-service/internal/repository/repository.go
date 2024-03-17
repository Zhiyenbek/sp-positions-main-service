package repository

import (
	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	AuthRepository
	TokenRepository
	RecruiterRepository
	CandidateRepository
}

type AuthRepository interface {
	GetUserInfoByLogin(login string) (string, string, error)
	Exists(loging string) (bool, error)
}

type RecruiterRepository interface {
	CreateRecruiter(input *models.RecruiterSignUpRequest) error
	Exists(publicID string) (bool, error)
}
type CandidateRepository interface {
	CreateCandidate(input *models.CandidateSignUpRequest) error
	Exists(publicID string) (bool, error)
}

type TokenRepository interface {
	SetRTToken(token *models.Token) error
	UnsetRTToken(publicID string) error
	GetToken(publicID string) (string, error)
}

func New(db *pgxpool.Pool, cfg *config.Configs, redis *redis.Client, log *zap.SugaredLogger) *Repository {
	return &Repository{
		AuthRepository:      NewAuthRepository(db, cfg.DB, log),
		TokenRepository:     NewTokenRepository(redis),
		RecruiterRepository: NewRecruiterRepository(db, cfg.DB, log),
		CandidateRepository: NewCandidateRepository(db, cfg.DB, log),
	}
}
