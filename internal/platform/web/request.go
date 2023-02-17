package web

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	defaultPageSize = 10
	defaultPage     = 1
)

func QueryPageParams(r *http.Request) (int, int, error) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	pageInt, pageSizeInt := defaultPage, defaultPageSize

	if len(page) > 0 {
		p, err := strconv.Atoi(page)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid page %s, must be a number", page)
		}
		pageInt = p
	}

	if len(pageSize) > 0 {
		ps, err := strconv.Atoi(pageSize)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid page size %s, must be a number", pageSize)
		}
		pageSizeInt = ps
	}

	return pageInt, pageSizeInt, nil
}
