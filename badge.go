package vrop

import (
	"fmt"
)

// Badge is a major or minor badge.
type Badge struct {
	// Type of the Badge
	Type string `json:"type,omitempty"`
	// Color of the Badge as determined by the system
	Color string `json:"color,omitempty"`
	// Score (value) associated with the Badge. This number represents the
	// absolute value of the Badge. Typically the value is between 0-100
	// but this is not the case all the time.
	Score float64 `json:"score,omitempty"`
}

func unpackBadge(m interface{}) (*Badge, error) {
	var pim map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		pim = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	p := &Badge{}
	for k, v := range pim {
		switch k {
		case "type":
			p.Type = v.(string)
		case "color":
			p.Color = v.(string)
		case "score":
			p.Score = v.(float64)
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return p, nil
}
