package response

import "github.com/google/uuid"

type BusinessResponse struct {
	Id           uuid.UUID             `json:"id"`
	Name         string                `json:"business_name"`
	OwnerName    string                `json:"owner_name"`
	Image        *string               `json:"image"`
	IsActive     bool                  `json:"is_active"`
	BusinessType *BusinessTypeResponse `json:"business_type"` // <- ubah jadi pointer
}
