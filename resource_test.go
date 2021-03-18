package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffsetPagination_GeneratePagination(t *testing.T) {
	var tests = map[string]struct{
		pagination OffsetPagination
		result Links
	}{
		"0 offset": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=0",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyNextPage: "/?page[limit]=100&page[offset]=100",
				KeyLastPage: "/?page[limit]=100&page[offset]=300",
			},

		},
		"0 offset and total a multiple of limit": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=0",
				Limit: 100,
				Total: 300,
			},
			result:     Links{
				KeyNextPage: "/?page[limit]=100&page[offset]=100",
				KeyLastPage: "/?page[limit]=100&page[offset]=200",
			},

		},
		"Offset below limit": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=80",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: "/?page[limit]=100&page[offset]=0",
				KeyNextPage: "/?page[limit]=100&page[offset]=180",
				KeyLastPage: "/?page[limit]=100&page[offset]=280",
			},

		},
		"Mid range offset": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=111",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: "/?page[limit]=100&page[offset]=0",
				KeyPreviousPage: "/?page[limit]=100&page[offset]=11",
				KeyNextPage: "/?page[limit]=100&page[offset]=211",
				KeyLastPage: "/?page[limit]=100&page[offset]=311",
			},

		},
		"Offset with other further params untouched": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=111&page[sort]=-1&aparam=2",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: "/?page[limit]=100&page[offset]=0&page[sort]=-1&aparam=2",
				KeyPreviousPage: "/?page[limit]=100&page[offset]=11&page[sort]=-1&aparam=2",
				KeyNextPage: "/?page[limit]=100&page[offset]=211&page[sort]=-1&aparam=2",
				KeyLastPage: "/?page[limit]=100&page[offset]=311&page[sort]=-1&aparam=2",
			},

		},
		"Offset with other previous params untouched": {
			pagination: OffsetPagination{
				URL:   "/?page[sort]=-1&aparam=2&page[limit]=111&page[offset]=111",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=0",
				KeyPreviousPage: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=11",
				KeyNextPage: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=211",
				KeyLastPage: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=311",
			},

		},
		"Offset with other params untouched": {
			pagination: OffsetPagination{
				URL:   "/?page[sort]=-1&page[limit]=111&aparam=2&page[offset]=111&lastparam=owt",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=0&lastparam=owt",
				KeyPreviousPage: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=11&lastparam=owt",
				KeyNextPage: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=211&lastparam=owt",
				KeyLastPage: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=311&lastparam=owt",
			},

		},
		"No params set": {
			pagination: OffsetPagination{
				URL:   "/",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyNextPage: "/?page[limit]=100&page[offset]=100",
				KeyLastPage: "/?page[limit]=100&page[offset]=300",
			},

		},
		"No paging set": {
			pagination: OffsetPagination{
				URL:   "/?param=owt",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyNextPage: "/?param=owt&page[limit]=100&page[offset]=100",
				KeyLastPage: "/?param=owt&page[limit]=100&page[offset]=300",
			},

		},
		"Non numeric parameter values": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=abc&page[offset]=def",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyLastPage: "/?page[limit]=100&page[offset]=300",
				KeyNextPage: "/?page[limit]=100&page[offset]=100",
			},

		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			underTest := test.pagination
			assert.Equal(t, test.result, *underTest.GeneratePagination())
		})
	}
}
