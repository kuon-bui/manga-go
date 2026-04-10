package common_test

import (
	"testing"

	"manga-go/internal/pkg/common"

	"github.com/stretchr/testify/assert"
)

func TestPaging_Fulfill_Defaults(t *testing.T) {
	p := &common.Paging{}
	p.Fulfill()

	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.Limit)
}

func TestPaging_Fulfill_NegativePage(t *testing.T) {
	p := &common.Paging{Page: -5, Limit: 10}
	p.Fulfill()

	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 10, p.Limit)
}

func TestPaging_Fulfill_NegativeLimit(t *testing.T) {
	p := &common.Paging{Page: 2, Limit: -1}
	p.Fulfill()

	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 20, p.Limit)
}

func TestPaging_Fulfill_ValidValues(t *testing.T) {
	p := &common.Paging{Page: 3, Limit: 50}
	p.Fulfill()

	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 50, p.Limit)
}

func TestPaging_GetLimit(t *testing.T) {
	p := &common.Paging{Page: 1, Limit: 15}

	assert.Equal(t, 15, p.GetLimit())
}

func TestPaging_GetLimit_Zero(t *testing.T) {
	p := &common.Paging{}

	assert.Equal(t, 20, p.GetLimit())
}

func TestPaging_GetOffset_FirstPage(t *testing.T) {
	p := &common.Paging{Page: 1, Limit: 20}

	assert.Equal(t, 0, p.GetOffset())
}

func TestPaging_GetOffset_SecondPage(t *testing.T) {
	p := &common.Paging{Page: 2, Limit: 20}

	assert.Equal(t, 20, p.GetOffset())
}

func TestPaging_GetOffset_ThirdPageCustomLimit(t *testing.T) {
	p := &common.Paging{Page: 3, Limit: 10}

	assert.Equal(t, 20, p.GetOffset())
}

func TestPaging_GetOffset_ZeroPage(t *testing.T) {
	// Zero page should be corrected to 1, so offset should be 0
	p := &common.Paging{Page: 0, Limit: 20}

	assert.Equal(t, 0, p.GetOffset())
}

func TestGenerateModelCode(t *testing.T) {
	tests := []struct {
		counter      int
		numberLength int
		prefix       string
		expected     string
	}{
		{1, 4, "AUT", "AUT0001"},
		{10, 4, "AUT", "AUT0010"},
		{100, 4, "AUT", "AUT0100"},
		{1000, 4, "AUT", "AUT1000"},
		{1, 6, "GEN", "GEN000001"},
		{0, 3, "TAG", "TAG000"},
	}

	for _, tt := range tests {
		result := common.GenerateModelCode(tt.counter, tt.numberLength, tt.prefix)
		assert.Equal(t, tt.expected, result)
	}
}
