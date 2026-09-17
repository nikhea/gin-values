package dto

// Pagination carries common page/page_size query params.
// Defaults: page=1, page_size=10 (capped at 100).
type Pagination struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// Normalize applies defaults and clamps to sane bounds.
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset returns the SQL offset for the current page.
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PageMeta is the pagination envelope returned with list endpoints.
type PageMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewPageMeta builds the envelope from a total row count.
func NewPageMeta(p Pagination, total int64) PageMeta {
	totalPages := 0
	if p.PageSize > 0 {
		totalPages = int((total + int64(p.PageSize) - 1) / int64(p.PageSize))
	}
	return PageMeta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
