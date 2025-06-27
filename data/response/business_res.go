package response

import "time"

type BusinessResponse struct {
	Id             int                      `json:"id"`
	Name           string                   `json:"business_name"`
	OwnerName      string                   `json:"owner_name"`
	BusinessTypeId int                      `json:"business_type_id"`
	Image          *string                  `json:"image,omitempty"`
	IsActive       bool                     `json:"is_active"`
	Branches       []BusinessBranchResponse `json:"branches,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}
