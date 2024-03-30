package models

type Position struct {
	PublicID          *string   `json:"public_id"`
	Name              *string   `json:"name"`
	Status            *int      `json:"status"`
	Skills            []*string `json:"skills"`
	Company           *Company  `json:"company,omitempty"`
	RecruiterPublicID *string   `json:"recruiter_public_id,omitempty"`
	Description       *string   `json:"description"`
}

type Company struct {
	PublicID    *string `json:"public_id"`
	Name        *string `json:"name"`
	Logo        *string `json:"logo"`
	Description *string `json:"description"`
}

type Question struct {
	ID               int    `json:"-"`
	PublicID         string `json:"public_id"`
	Name             string `json:"name"`
	PositionPublicID string `json:"-"`
	PositionID       int    `json:"-"`
	ReadDuration     int    `json:"read_duration"`
	AnswerDuration   int    `json:"answer_duration"`
}
