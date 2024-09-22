package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPagedResponse(t *testing.T) {
	testCases := []struct {
		name               string
		numberOfItems      int
		page               int
		pageSize           int
		expectedFirstItem  int
		expectedLastItem   int
		expectedTotalPages int
	}{
		{
			name:               "test 1",
			numberOfItems:      10,
			page:               1,
			pageSize:           5,
			expectedFirstItem:  1,
			expectedLastItem:   5,
			expectedTotalPages: 2,
		},
		{
			name:               "test 2",
			numberOfItems:      2,
			page:               1,
			pageSize:           5,
			expectedFirstItem:  1,
			expectedLastItem:   2,
			expectedTotalPages: 1,
		},
		{
			name:               "test 3",
			numberOfItems:      9,
			page:               2,
			pageSize:           5,
			expectedFirstItem:  6,
			expectedLastItem:   9,
			expectedTotalPages: 2,
		},
		{
			name:               "test 4",
			numberOfItems:      10,
			page:               0,
			pageSize:           0,
			expectedFirstItem:  1,
			expectedLastItem:   1,
			expectedTotalPages: 10,
		},
		{
			name:               "test 5",
			numberOfItems:      0,
			page:               1,
			pageSize:           5,
			expectedFirstItem:  1,
			expectedLastItem:   5,
			expectedTotalPages: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var items []int
			for i := 1; i <= tc.numberOfItems; i++ {
				items = append(items, i)
			}

			e := echo.New()
			url := fmt.Sprintf("/test?page=%d&page_size=%d", tc.page, tc.pageSize)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			c := e.NewContext(req, httptest.NewRecorder())

			pagedResponse := NewPagedResponse(c, items)
			items, ok := pagedResponse.Items.([]int)
			require.True(t, ok)

			if tc.numberOfItems == 0 {
				require.Equal(t, 0, len(items))
			} else {
				require.Equal(t, (tc.expectedLastItem-tc.expectedFirstItem)+1, len(items))
				require.Equal(t, tc.expectedFirstItem, items[0])
				require.Equal(t, tc.expectedLastItem, items[len(items)-1])
			}

			require.Equal(t, max(tc.page, 1), pagedResponse.Metadata.Page)
			require.Equal(t, max(tc.pageSize, 1), pagedResponse.Metadata.PageSize)
			require.Equal(t, tc.numberOfItems, pagedResponse.Metadata.TotalItems)
		})
	}
}
