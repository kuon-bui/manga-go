package common

type Paging struct {
	Limit int `json:"limit" form:"limit"`
	Page  int `json:"page" form:"page"`
}

func (p *Paging) Fulfill() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 20
	}
}

func (p *Paging) GetLimit() int {
	p.Fulfill()
	return p.Limit
}

func (p *Paging) GetOffset() int {
	p.Fulfill()
	return (p.Page - 1) * p.Limit
}
