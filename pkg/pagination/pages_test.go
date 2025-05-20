package pagination

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                 string
		page, perPage, total int
		wantPage             int
		wantPerPage          int
		wantTotal            int
		wantPageCount        int
		wantOffset           int
		wantLimit            int
	}{
		{"normal", 2, 20, 50, 2, 20, 50, 3, 20, 20},
		{"too big page", 10, 20, 50, 3, 20, 50, 3, 40, 20},
		{"zero page", 0, 20, 50, 1, 20, 50, 3, 0, 20},
		{"negative perPage", 1, -1, 50, 1, 100, 50, 1, 0, 100},
		{"too big perPage", 1, 10000, 50, 1, 1000, 50, 1, 0, 1000},
		{"zero total", 1, 20, 0, 1, 20, 0, 0, 0, 20},
		{"negative total", 1, 20, -1, 1, 20, -1, -1, 0, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.page, tt.perPage, tt.total)
			assert.Equal(t, tt.wantPage, p.Page)
			assert.Equal(t, tt.wantPerPage, p.PerPage)
			assert.Equal(t, tt.wantTotal, p.TotalCount)
			assert.Equal(t, tt.wantPageCount, p.PageCount)
			assert.Equal(t, tt.wantOffset, p.Offset())
			assert.Equal(t, tt.wantLimit, p.Limit())
		})
	}
}

func TestNewFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com?page=2&per_page=30", bytes.NewBuffer(nil))
	p := NewFromRequest(req, 90)

	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 30, p.PerPage)
	assert.Equal(t, 90, p.TotalCount)
	assert.Equal(t, 3, p.PageCount)
	assert.Equal(t, 30, p.Offset())
	assert.Equal(t, 30, p.Limit())
}

func Test_parseInt(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		defValue int
		want     int
	}{
		{"valid", "42", 10, 42},
		{"empty", "", 10, 10},
		{"invalid", "abc", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseInt(tt.value, tt.defValue))
		})
	}
}
