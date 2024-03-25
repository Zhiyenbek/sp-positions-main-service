package service

import (
	"github.com/Zhiyenbek/users-auth-service/config"
	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/Zhiyenbek/users-auth-service/internal/repository"
	"go.uber.org/zap"
)

type AuthService interface {
	CandidateLogin(creds *models.UserSignInRequest) (*models.Tokens, error)
	RecruiterLogin(creds *models.UserSignInRequest) (*models.Tokens, error)
	RefreshToken(tokenString string) (*models.Tokens, error)
	CreateRecruiter(req *models.RecruiterSignUpRequest) error
	CreateCandidate(req *models.CandidateSignUpRequest) error
	SignOut(accessToken string) error
}

type Service struct {
	AuthService
}

func New(repos *repository.Repository, log *zap.SugaredLogger, cfg *config.Configs) *Service {
	return &Service{
		AuthService: NewAuthService(repos, cfg, log),
	}
}
