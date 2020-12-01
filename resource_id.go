package vrop

import (
	"fmt"
)

// ResourceIdentifier is a reference to an object.
type ResourceIdentifier struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func unpackResourceIdentifier(m interface{}) (*ResourceIdentifier, error) {
	var pm map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		pm = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	p := &ResourceIdentifier{}

	var name, dataType, value string
	// var isPartOfUniqueness bool
	if pv, exists := pm["identifierType"]; exists {
		it := pv.(map[string]interface{})
		for k, v := range it {
			switch k {
			case "name":
				name = v.(string)
			case "dataType":
				dataType = v.(string)
			case "isPartOfUniqueness":
				// isPartOfUniqueness = v.(bool)
			default:
				return nil, fmt.Errorf("resource id contains unsupported identifierType: %s, data: %v", k, m)
			}
		}
	}

	if pv, exists := pm["value"]; exists {
		switch dataType {
		case "STRING":
			value = pv.(string)
		default:
			return nil, fmt.Errorf("resource id contains unsupported type: %s, data: %v", dataType, m)
		}
	}

	if name == "" {
		return nil, fmt.Errorf("resource id name not found, data: %v", m)
	}
	//if value == "" {
	//	return nil, fmt.Errorf("resource id value not found, data: %v", m)
	//}

	p.Key = name
	p.Value = value

	return p, nil
}
