package jsonapi

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Payloader is used to encapsulate the One and Many payload types
type Payloader interface {
	clearIncluded()
	AddPagination(paginator Paginator)
}

// NulledPayload allows for raw message to inspect nulls
type NulledPayload struct {
	Data ResourceObjNulls `json:"data"`
}

// OnePayload is used to represent a generic JSON API payload where a single
// resource (ResourceObj) was included as an {} in the "data" key
type OnePayload struct {
	Data     *ResourceObj   `json:"data"`
	Included []*ResourceObj `json:"included,omitempty"`
	Links    *Links         `json:"links,omitempty"`
	Meta     *Meta          `json:"meta,omitempty"`
}

func (p *OnePayload) clearIncluded() {
	p.Included = []*ResourceObj{}
}

func (p *OnePayload) AddPagination(paginator Paginator) {

}

// ManyPayload is used to represent a generic JSON API payload where many
// resources (Nodes) were included in an [] in the "data" key
type ManyPayload struct {
	Data     []*ResourceObj `json:"data"`
	Included []*ResourceObj `json:"included,omitempty"`
	Links    *Links         `json:"links,omitempty"`
	Meta     *Meta          `json:"meta,omitempty"`
}

func (p *ManyPayload) clearIncluded() {
	p.Included = []*ResourceObj{}
}

func (p *ManyPayload) AddPagination(paginator Paginator) {
	p.Links = paginator.GeneratePagination()
}

// ResourceObjNulls is used to represent a generic JSON API Resource with null fields
type ResourceObjNulls struct {
	Type       string                     `json:"type"`
	ID         string                     `json:"id,omitempty"`
	Attributes map[string]json.RawMessage `json:"attributes,omitempty"`
}

// ResourceObj is used to represent a generic JSON API Resource
type ResourceObj struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
	Links         *Links                 `json:"links,omitempty"`
	Meta          *Meta                  `json:"meta,omitempty"`
}

// RelationshipOneNode is used to represent a generic has one JSON API relation
type RelationshipOneNode struct {
	Data  *ResourceObj `json:"data"`
	Links *Links       `json:"links,omitempty"`
	Meta  *Meta        `json:"meta,omitempty"`
}

// RelationshipManyNode is used to represent a generic has many JSON API
// relation
type RelationshipManyNode struct {
	Data  []*ResourceObj `json:"data"`
	Links *Links         `json:"links,omitempty"`
	Meta  *Meta          `json:"meta,omitempty"`
}

// Links is used to represent a `links` object.
// http://jsonapi.org/format/#document-links
type Links map[string]interface{}

func (l *Links) validate() (err error) {
	// Each member of a links object is a “link”. A link MUST be represented as
	// either:
	//  - a string containing the link’s URL.
	//  - an object (“link object”) which can contain the following members:
	//    - href: a string containing the link’s URL.
	//    - meta: a meta object containing non-standard meta-information about the
	//            link.
	for k, v := range *l {
		_, isString := v.(string)
		_, isLink := v.(Link)

		if !(isString || isLink) {
			return fmt.Errorf(
				"The %s member of the links object was not a string or link object",
				k,
			)
		}
	}
	return
}

// Link is used to represent a member of the `links` object.
type Link struct {
	Href string `json:"href"`
	Meta Meta   `json:"meta,omitempty"`
}

// Linkable is used to include document links in response data
// e.g. {"self": "http://example.com/posts/1"}
type Linkable interface {
	JSONAPILinks() *Links
}

// RelationshipLinkable is used to include relationship links  in response data
// e.g. {"related": "http://example.com/posts/1/comments"}
type RelationshipLinkable interface {
	// JSONAPIRelationshipLinks will be invoked for each relationship with the corresponding relation name (e.g. `comments`)
	JSONAPIRelationshipLinks(relation string) *Links
}

// Meta is used to represent a `meta` object.
// http://jsonapi.org/format/#document-meta
type Meta map[string]interface{}

// Metable is used to include document meta in response data
// e.g. {"foo": "bar"}
type Metable interface {
	JSONAPIMeta() *Meta
}

// RelationshipMetable is used to include relationship meta in response data
type RelationshipMetable interface {
	// JSONRelationshipMeta will be invoked for each relationship with the corresponding relation name (e.g. `comments`)
	JSONAPIRelationshipMeta(relation string) *Meta
}

// Paginator allows clients to paginate the result set
type Paginator interface {
	GeneratePagination() *Links
}

type OffsetPagination struct {
	URL   string
	Limit int64
	Total int64
}

func (p *OffsetPagination) GeneratePagination() *Links {
	if p.Total < p.Limit { // no pagination needed
		return nil
	}

	// initiate the URL - if the page offset and Limit have not been set or is devoid of all
	// query parameters then initialising will make string replacement a simple operation

	if !strings.Contains(p.URL, "page[limit]") {
		p.appendToURL("page[limit]=" + strconv.FormatInt(p.Limit, 10))
	}
	if !strings.Contains(p.URL, "page[offset]") {
		p.appendToURL("page[offset]=0")
	}

	links := Links{}
	limit := int64(math.Min(float64(getPageParam("Limit", p.URL)), float64(p.Limit)))
	if limit == 0 {
		limit = p.Limit
	}
	offset := int64(math.Max(float64(getPageParam("offset", p.URL)), float64(0)))

	if offset > 0 {
		firstUrl := p.URL
		replaceParam(&firstUrl, `page[limit]`, strconv.FormatInt(limit, 10))
		replaceParam(&firstUrl, `page[offset]`, strconv.FormatInt(0, 10))
		links[KeyFirstPage] = firstUrl
	}

	if offset > limit {
		prevUrl := p.URL
		replaceParam(&prevUrl, `page[limit]`, strconv.FormatInt(limit, 10))
		prevOffset := offset-limit
		replaceParam(&prevUrl, `page[offset]`, strconv.FormatInt(prevOffset, 10))
		links[KeyPreviousPage] = prevUrl
	}

	if offset+limit < p.Total-limit {
		nextUrl := p.URL
		replaceParam(&nextUrl, `page[limit]`, strconv.FormatInt(limit, 10))
		nextOffset := offset + limit
		replaceParam(&nextUrl, `page[offset]`, strconv.FormatInt(nextOffset, 10))
		links[KeyNextPage] = nextUrl
	}

	if offset+limit < p.Total {
		lastUrl := p.URL
		replaceParam(&lastUrl, `page[limit]`, strconv.FormatInt(limit, 10))
		pages := p.Total / limit
		if p.Total%limit > 0 {
			pages += 1
		}
		lastOffset := ((pages-1)*limit)
		offsetShift := offset % limit
		lastOffset += offsetShift
		if lastOffset > p.Total {
			lastOffset -= limit
		}
		replaceParam(&lastUrl, `page[offset]`, strconv.FormatInt(lastOffset, 10))
		links[KeyLastPage] = lastUrl
	}

	return &links
}

func getPageParam(name, url string) int64 {
	val := 0
	valRe := regexp.MustCompile(fmt.Sprintf(`page\[%s\]=(\d+)`, name))
	match := valRe.FindStringSubmatch(url)
	if len(match) == 2 { // when we have found the \d portion
		ql := match[1]
		val, _ = strconv.Atoi(ql)
	}
	return int64(val)
}

func replaceParam(url *string, param, value string) {
	var sb strings.Builder
	sb.WriteString(param)
	sb.WriteString("=")
	sb.WriteString(value)
	newParam := sb.String()

	seek := fmt.Sprintf(`%s=[^&]+`, regexSafe(param))
	regex := regexp.MustCompile(seek)
	match := regex.ReplaceAllString(*url, newParam)

	*url = match
}

func regexSafe(in string) string {
	chars := []string{"]", "^", "\\", "[", ".", "(", ")", "-"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("([" + r + "])+")
	out := re.ReplaceAllString(in, "\\$1")
	return out
}

func (p *OffsetPagination) appendToURL(param string) {
	if !strings.Contains(p.URL, "?") {
		p.URL += "?" + param
	} else {
		p.URL += "&" + param
	}
}
