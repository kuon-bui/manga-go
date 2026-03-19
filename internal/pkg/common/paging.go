package common

import "strings"

type Paging struct {
	Limit int   `json:"limit" form:"limit"`
	Page  int   `json:"page" form:"page"`
	Total int64 `json:"total" form:"total"`
	//	Support cursor with UID
	FakeCursor string `json:"cursor" form:"cursor"`
	NextCursor string `json:"next_cursor" form:"next_cursor"`
}

func (p *Paging) Fulfill() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit == 0 {
		p.Limit = 20
	}

	p.FakeCursor = strings.TrimSpace(p.FakeCursor)
}

func (p *Paging) GetLimit() int {
	p.Fulfill()
	return p.Limit
}

func (p *Paging) GetOffset() int {
	p.Fulfill()
	return (p.Page - 1) * p.Limit
}
