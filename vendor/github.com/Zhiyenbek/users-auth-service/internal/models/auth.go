package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type UserSignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserData struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RecruiterSignUpRequest struct {
	UserData
	CompanyName     string    `json:"company_name"`
	CompanyPublicID uuid.UUID `json:"company_public_id"`
}

type CandidateSignUpRequest struct {
	UserData
	Resume          string   `json:"resume"`
	CurrentPosition string   `json:"current_position"`
	Bio             string   `json:"bio"`
	Education       string   `json:"education"`
	Skills          []string `json:"skills"`
}

type Token struct {
	PublicID   string
	TokenValue string
	Role       string
	TTL        time.Duration
}

// Tokens - structure for holding access and refresh token
type Tokens struct {
	AccessToken  *Token
	RefreshToken *Token
}
type JwtUserClaims struct {
	PublicID string `json:"user_public_id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
