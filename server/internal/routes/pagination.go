package routes

import (
	"github.com/labstack/echo/v4"
	"reflect"
	"strconv"
)

type pageMetadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type PagedResponse struct {
	Metadata pageMetadata `json:"metadata"`
	Items    any          `json:"items"`
}

func NewPagedResponse(c echo.Context, items any) PagedResponse {
	v := reflect.ValueOf(items)

	if v.Kind() != reflect.Slice {
		panic("expected items to be a slice")
	}

	page, pageSize := getPagePrams(c)
	totalItems := v.Len()
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	return PagedResponse{
		Metadata: pageMetadata{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: (totalItems + pageSize - 1) / pageSize,
		},
		Items: v.Slice(start, end).Interface(),
	}
}

func getPagePrams(c echo.Context) (int, int) {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 1
	}
	return page, pageSize
}
