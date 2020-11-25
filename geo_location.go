package vrop

import (
	"fmt"
)

// GeoLocation is geographical location.
type GeoLocation struct {
	// Latitude of the location.
	Latitude float64 `json:"latitude,omitempty"`
	// Longitude of the location.
	Longitude float64 `json:"longitude,omitempty"`
}

func unpackGeoLocation(m interface{}) (*GeoLocation, error) {
	var pim map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		pim = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	p := &GeoLocation{}
	for k, v := range pim {
		switch k {
		case "latitude":
			p.Latitude = v.(float64)
		case "longitude":
			p.Longitude = v.(float64)
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return p, nil
}
