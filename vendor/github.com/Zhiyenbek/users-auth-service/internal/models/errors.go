package models

import "errors"

var (
	ErrInvalidInput          = errors.New("INVALID_INPUT")
	ErrInternalServer        = errors.New("INTERNAL_SERVER_ERROR")
	ErrInvalidPasswordFormat = errors.New("WRONG_PASSWORD_FORMAT")
	ErrWrongCredential       = errors.New("WRONG_CREDENTIALS")
	ErrInvalidToken          = errors.New("INVALID_TOKEN")
	ErrTokenExpired          = errors.New("TOKEN_EXPIRED")
	ErrCompanyDoesntExists   = errors.New("COMPANY_DOES_NOT_EXIST")
	ErrUsernameExists        = errors.New("USERNAME_EXISTS")
)
