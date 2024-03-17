package models

type Position struct {
	PublicID          *string   `json:"public_id"`
	Name              *string   `json:"name"`
	Status            *int      `json:"status"`
	Skills            []*string `json:"skills"`
	Company           *Company  `json:"company"`
	RecruiterPublicID *string   `json:"recruiter_public_id,omitempty"`
	Description       *string   `json:"description,omitempty"`
}

type Company struct {
	PublicID    *string `json:"public_id"`
	Name        *string `json:"name"`
	Logo        *string `json:"logo,omitempty"`
	Description *string `json:"description,omitempty"`
}
