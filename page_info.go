package vrop

import (
	"fmt"
)

// PageInfo is paging reference.
type PageInfo struct {
	// Total number of results.
	Total int `json:"totalCount,omitempty"`
	// The current page number.
	ID int `json:"page,omitempty"`
	// Number of entries allowed in a page.
	Size int `json:"pageSize,omitempty"`
	// A CSV list of field names. Usually found in an enumeration, for example
	SortBy string `json:"sortBy,omitempty"`
	// A CSV list of values. If not specified or if list shorter than sortFields then SortOrder.ASCENDING is assumed.
	SortOrder string `json:"sortOrder,omitempty"`
}

func unpackPageInfo(m interface{}) (*PageInfo, error) {
	var pim map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		pim = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	p := &PageInfo{}
	for k, v := range pim {
		switch k {
		case "page":
			p.ID = int(v.(float64))
		case "pageSize":
			p.Size = int(v.(float64))
		case "totalCount":
			p.Total = int(v.(float64))
		case "sortBy":
			p.SortBy = v.(string)
		case "sortOrder":
			p.SortOrder = v.(string)
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return p, nil
}
