package requests

import "github.com/grassbusinesslabs/eventio-go-back/internal/domain"

type Page struct {
	Page int64 `json:"page" validate:"required"`
}

const pageSize = 15

func (r Page) ToDomainModel() (interface{}, error) {
	return domain.Pagination{
		Page:         uint64(r.Page),
		CountPerPage: pageSize,
	}, nil
}
