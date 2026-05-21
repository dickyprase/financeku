package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

type Params struct {
	Page    int
	PerPage int
}

func Parse(r *http.Request) Params {
	page := parseIntParam(r, "page", DefaultPage)
	perPage := parseIntParam(r, "per_page", DefaultPerPage)

	if page < 1 {
		page = DefaultPage
	}
	if perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	return Params{
		Page:    page,
		PerPage: perPage,
	}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func TotalPages(total int64, perPage int) int {
	pages := int(total) / perPage
	if int(total)%perPage > 0 {
		pages++
	}
	return pages
}

func parseIntParam(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
