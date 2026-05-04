package common

import "testing"

func TestPagingFulfillDefaults(t *testing.T) {
	p := &Paging{}

	p.Fulfill()

	if p.Page != 1 {
		t.Fatalf("expected default page = 1, got %d", p.Page)
	}
	if p.Limit != 20 {
		t.Fatalf("expected default limit = 20, got %d", p.Limit)
	}
}

func TestPagingGetLimitAppliesDefaults(t *testing.T) {
	p := &Paging{Limit: 0, Page: -1}

	got := p.GetLimit()

	if got != 20 {
		t.Fatalf("expected GetLimit() = 20, got %d", got)
	}
	if p.Page != 1 {
		t.Fatalf("expected page to be normalized to 1, got %d", p.Page)
	}
}

func TestPagingGetOffset(t *testing.T) {
	p := &Paging{Page: 3, Limit: 15}

	got := p.GetOffset()

	if got != 30 {
		t.Fatalf("expected offset = 30, got %d", got)
	}
}
