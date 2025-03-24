package dto

import "strings"

type MetaPagination struct {
	Page      int    `json:"page" query:"page"`
	Limit     int    `json:"limit" query:"limit"`
	Order     string `json:"order,omitempty" query:"order"`
	SortBy    string `json:"sort_by,omitempty" query:"sort_by"`
	Offset    int    `json:"offset,omitempty"`
	Total     int64  `json:"total,omitempty"`
	TotalPage int64  `json:"total_page,omitempty"`
	BaseResponse
}

func (p *MetaPagination) ParsePagination() *MetaPagination {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit > 50 {
		p.Limit = 10
	}

	offset := (p.Page - 1) * p.Limit
	if offset < 0 {
		offset = 0
	}

	if p.SortBy == "" {
		p.SortBy = "created_at"
	}

	if p.Order == "" {
		p.Order = " DESC"
	} else {
		if strings.ToLower(p.Order) != "asc" && strings.ToLower(p.Order) != "desc" {
			p.Order = " DESC"
		}
	}

	p.Offset = offset
	return p
}
