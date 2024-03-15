package models

type Position struct {
	PublicID string `json:"public_id"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
}
