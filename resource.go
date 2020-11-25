package vrop

import (
	"fmt"
	"time"
)

// Resource is a resource record.
type Resource struct {
	// Identifier of the Resource (typically a UUID).
	ID string `json:"identifier,omitempty"`
	// Description of the Resource
	Description string `json:"description,omitempty"`
	// Time the Resource was created in the system.
	CreationTime time.Time `json:"creationTime,omitempty"`
	// Resource key of the Resource.
	Key *ResourceKey `json:"resourceKey,omitempty"`
	// Credential instance identifier assigned to this Resource
	CredentialInstanceID string `json:"credentialInstanceId,omitempty"`
	// Geographical location of the resource.
	GeoLocation *GeoLocation `json:"geoLocation,omitempty"`
	// The resource status and resource state for this resource
	// as reported by one or more adapter instances
	StatusStates []*ResourceStatusState `json:"resourceStatusStates,omitempty"`
	// Health of the Resource.
	Health string `json:"resourceHealth,omitempty"`
	// Resource Health Score.
	HealthValue float64 `json:"resourceHealthValue,omitempty"`
	// DT calculation enabled or not. By default DT calculation for
	// a resource is enabled (during its creation).
	DynamicThresholdEnabled bool `json:"dtEnabled,omitempty"`
	// The Monitoring Interval (in minutes) the Resource was configured with
	MonitoringInterval float64 `json:"monitoringInterval,omitempty"`
	// The various major and minor badges and their values for a Resource.
	Badges []*Badge `json:"badges,omitempty"`
	// Collection of related resource identifiers.
	RelatedResources []interface{} `json:"relatedResources,omitempty"`
	// Extension values that were added to the given object by third-party.
	Extension interface{} `json:"extension,omitempty"`
	// Set of useful links related to the current object.
	Links []*Link `json:"links,omitempty"`
}

func unpackResource(m interface{}) (*Resource, error) {
	var rmap map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		rmap = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	r := &Resource{}
	for k, v := range rmap {
		switch k {
		case "identifier":
			r.ID = v.(string)
		case "description":
			r.Description = v.(string)
		case "creationTime":
			secs := int64(v.(float64) / 1000)
			nsecs := int64(((v.(float64) / 1000) - float64(secs)) * 1e9)
			r.CreationTime = time.Unix(secs, nsecs)
		case "resourceKey":
			if s, err := unpackResourceKey(v.(interface{})); err != nil {
				return nil, fmt.Errorf("failed to unpack %s resourceKey: %s", k, err)
			} else {
				r.Key = s
			}
		case "credentialInstanceId":
			r.CredentialInstanceID = v.(string)
		case "geoLocation":
			if s, err := unpackGeoLocation(v.(interface{})); err != nil {
				return nil, fmt.Errorf("failed to unpack %s geoLocation: %s", k, err)
			} else {
				r.GeoLocation = s
			}
		case "resourceStatusStates":
			for _, item := range v.([]interface{}) {
				if s, err := unpackResourceStatusState(item); err != nil {
					return nil, fmt.Errorf("failed to unpack %s resourceStatusStates: %s", k, err)
				} else {
					r.StatusStates = append(r.StatusStates, s)
				}
			}
		case "resourceHealth":
			r.Health = v.(string)
		case "resourceHealthValue":
			r.HealthValue = v.(float64)
		case "dtEnabled":
			r.DynamicThresholdEnabled = v.(bool)
		case "monitoringInterval":
			r.MonitoringInterval = v.(float64)
		case "badges":
			for _, item := range v.([]interface{}) {
				if badge, err := unpackBadge(item); err != nil {
					return nil, fmt.Errorf("failed to unpack %s badge: %s", k, err)
				} else {
					r.Badges = append(r.Badges, badge)
				}
			}
		case "relatedResources":
			r.RelatedResources = v.([]interface{})
		case "extension":
			r.Extension = v.(interface{})
		case "links":
			for _, item := range v.([]interface{}) {
				if link, err := unpackLink(item); err != nil {
					return nil, fmt.Errorf("failed to unpack %s link: %s", k, err)
				} else {
					r.Links = append(r.Links, link)
				}
			}
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return r, nil
}
