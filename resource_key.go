package vrop

import (
	"fmt"
)

// ResourceKey is a key of the Resource.
type ResourceKey struct {
	// Name of the Resource
	Name string `json:"name,omitempty"`
	// Adapter Kind to which the resource belongs to
	AdapterKindKey string `json:"adapterKindKey,omitempty"`
	// Resource Kind to which the resource belongs to
	ResourceKindKey string `json:"resourceKindKey,omitempty"`
	// A collection of Resource Identifiers. A Resource Identifier is a
	// key-value pair that encapsulates the identity of the resource
	ResourceIdentifiers []*ResourceIdentifier `json:"resourceIdentifiers,omitempty"`
	// Set of useful links related to the current object.
	Links []*Link `json:"links,omitempty"`
	// Extension values that were added to the given object by third-party.
	Extension interface{} `json:"extension,omitempty"`
}

func unpackResourceKey(m interface{}) (*ResourceKey, error) {
	var pim map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		pim = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	p := &ResourceKey{}
	for k, v := range pim {
		switch k {
		case "name":
			p.Name = v.(string)
		case "adapterKindKey":
			p.AdapterKindKey = v.(string)
		case "resourceKindKey":
			p.ResourceKindKey = v.(string)
		case "resourceIdentifiers":
			for _, item := range v.([]interface{}) {
				resourceID, err := unpackResourceIdentifier(item)
				if err != nil {
					return nil, fmt.Errorf("failed to unpack %s resourceIdentifier: %s", k, err)
				}
				p.ResourceIdentifiers = append(p.ResourceIdentifiers, resourceID)
			}
		case "links":
			for _, item := range v.([]interface{}) {
				link, err := unpackLink(item)
				if err != nil {
					return nil, fmt.Errorf("failed to unpack %s link: %s", k, err)
				}
				p.Links = append(p.Links, link)
			}
		case "extension":
			// TODO
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return p, nil
}
