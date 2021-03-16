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
		"Offset": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=111",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: Link{Href: "/?page[limit]=100&page[offset]=0"},
				KeyPreviousPage: Link{Href: "/?page[limit]=100&page[offset]=11"},
				KeyNextPage: Link{Href: "/?page[limit]=100&page[offset]=211"},
				KeyLastPage: Link{Href: "/?page[limit]=100&page[offset]=311"},
			},

		},
		"Offset with other further params untouched": {
			pagination: OffsetPagination{
				URL:   "/?page[limit]=111&page[offset]=111&page[sort]=-1&aparam=2",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: Link{Href: "/?page[limit]=100&page[offset]=0&page[sort]=-1&aparam=2"},
				KeyPreviousPage: Link{Href: "/?page[limit]=100&page[offset]=11&page[sort]=-1&aparam=2"},
				KeyNextPage: Link{Href: "/?page[limit]=100&page[offset]=211&page[sort]=-1&aparam=2"},
				KeyLastPage: Link{Href: "/?page[limit]=100&page[offset]=311&page[sort]=-1&aparam=2"},
			},

		},
		"Offset with other previous params untouched": {
			pagination: OffsetPagination{
				URL:   "/?page[sort]=-1&aparam=2&page[limit]=111&page[offset]=111",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: Link{Href: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=0"},
				KeyPreviousPage: Link{Href: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=11"},
				KeyNextPage: Link{Href: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=211"},
				KeyLastPage: Link{Href: "/?page[sort]=-1&aparam=2&page[limit]=100&page[offset]=311"},
			},

		},
		"Offset with other params untouched": {
			pagination: OffsetPagination{
				URL:   "/?page[sort]=-1&page[limit]=111&aparam=2&page[offset]=111&lastparam=owt",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyFirstPage: Link{Href: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=0&lastparam=owt"},
				KeyPreviousPage: Link{Href: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=11&lastparam=owt"},
				KeyNextPage: Link{Href: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=211&lastparam=owt"},
				KeyLastPage: Link{Href: "/?page[sort]=-1&page[limit]=100&aparam=2&page[offset]=311&lastparam=owt"},
			},

		},
		"No params set": {
			pagination: OffsetPagination{
				URL:   "/",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyNextPage: Link{Href: "/?page[limit]=100&page[offset]=100"},
				KeyLastPage: Link{Href: "/?page[limit]=100&page[offset]=300"},
			},

		},
		"No paging set": {
			pagination: OffsetPagination{
				URL:   "/?param=owt",
				Limit: 100,
				Total: 334,
			},
			result:     Links{
				KeyNextPage: Link{Href: "/?param=owt&page[limit]=100&page[offset]=100"},
				KeyLastPage: Link{Href: "/?param=owt&page[limit]=100&page[offset]=300"},
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
