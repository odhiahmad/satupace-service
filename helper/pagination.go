package helper

import (
	"strings"

	"github.com/odhiahmad/kasirku-service/data/request"
	"gorm.io/gorm"
)

type Paginator struct {
	pagination request.Pagination
	validSort  map[string]bool
}

// Buat paginator baru dengan validasi field sort yang diizinkan
func Paginate(pagination request.Pagination) *Paginator {
	return &Paginator{
		pagination: pagination,
		validSort: map[string]bool{
			"name":        true,
			"created_at":  true,
			"updated_at":  true,
			"description": true,
		},
	}
}

// Apply pagination ke query GORM
func (p *Paginator) Paginate(db *gorm.DB, result interface{}) (int, int, error) {
	// Validasi nilai
	if p.pagination.Page <= 0 {
		p.pagination.Page = 1
	}
	if p.pagination.Limit <= 0 {
		p.pagination.Limit = 10
	}

	sortBy := p.pagination.SortBy
	if sortBy == "" || !p.validSort[sortBy] {
		sortBy = "created_at"
	}

	order := "asc"
	if strings.ToLower(p.pagination.OrderBy) == "desc" {
		order = "desc"
	}

	offset := (p.pagination.Page - 1) * p.pagination.Limit

	query := db.
		Order(sortBy + " " + order).
		Limit(p.pagination.Limit).
		Offset(offset)

	err := query.Find(result).Error
	if err != nil {
		return 0, 0, err
	}

	return p.pagination.Page, p.pagination.Limit, nil
}
