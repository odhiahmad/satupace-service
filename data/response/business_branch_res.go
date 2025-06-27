package response

import "time"

type BusinessBranchResponse struct {
	Id          int       `json:"id"`
	BusinessId  int       `json:"business_id"`
	PhoneNumber string    `json:"phone_number"`
	Rating      string    `json:"rating"`
	Provinsi    string    `json:"provinsi"`
	Kota        string    `json:"kota"`
	Kecamatan   string    `json:"kecamatan"`
	PostalCode  string    `json:"postal_code"`
	IsMain      bool      `json:"is_main"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
