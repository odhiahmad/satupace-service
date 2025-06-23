package helper

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

func Paginate(pagination request.Pagination) *paginator.Paginator {
	// Default values
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.SortBy == "" {
		pagination.SortBy = "id"
	}
	order := paginator.ASC
	if pagination.OrderBy == "desc" {
		order = paginator.DESC
	}

	opts := []paginator.Option{
		&paginator.Config{
			Keys:  []string{pagination.SortBy},
			Limit: pagination.Limit,
			Order: order,
		},
	}

	// Cursors
	if pagination.After != "" {
		opts = append(opts, paginator.WithAfter(pagination.After))
	}
	if pagination.Before != "" {
		opts = append(opts, paginator.WithBefore(pagination.Before))
	}

	return paginator.New(opts...)
}
