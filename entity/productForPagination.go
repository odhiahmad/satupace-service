package entity

import "time"

type ProductForPagination struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	BasePrice float64   `json:"base_price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductForPagination) TableName() string {
	return "products"
}
