package helper

import (
	"strings"

	"loka-kasir/data/request"

	"gorm.io/gorm"
)

type Paginator struct {
	pagination request.Pagination
	validSort  map[string]bool
}

func Paginate(pagination request.Pagination, validSortFields []string) *Paginator {
	validSort := make(map[string]bool)
	for _, field := range validSortFields {
		validSort[field] = true
	}

	return &Paginator{
		pagination: pagination,
		validSort:  validSort,
	}
}

func (p *Paginator) Paginate(db *gorm.DB, result interface{}) (int, int, error) {
	if p.pagination.Page <= 0 {
		p.pagination.Page = 1
	}
	if p.pagination.Limit <= 0 {
		p.pagination.Limit = 10
	}

	sortBy := strings.ToLower(p.pagination.SortBy)
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
