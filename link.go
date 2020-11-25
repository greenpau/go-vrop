package vrop

import (
	"fmt"
)

// Link is a reference to an object.
type Link struct {
	// 	The url part of this tag. Note that the href can be absolute or relative to the current url.
	// If the href begins with '/', then its absolute.
	// If the href begins with the protocol element(http, https), its absolute
	// Otherwise its relative.
	Href string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
	// NEXT: Used to represent the next page accessible using this url
	// PREVIOUS: Used to represent the previous page accessible using this url
	// START: Used to represent the first page accessible using this url
	// END: Used to represent the last page accessible using this url
	// RELATED: Used to represent that this link points to an object related to the link's parent
	// SELF: Used to represent that this link points to more information of the link's parent tag/object
	Relation string `json:"rel,omitempty"`
}

func unpackLink(m interface{}) (*Link, error) {
	var pim map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		pim = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	p := &Link{}
	for k, v := range pim {
		switch k {
		case "href":
			p.Href = v.(string)
		case "name":
			p.Name = v.(string)
		case "rel":
			p.Relation = v.(string)
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return p, nil
}
