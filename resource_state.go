package vrop

import (
	"fmt"
)

// ResourceStatusState is paging reference.
type ResourceStatusState struct {
	// The Adapter Instance identifier associated with the status and state
	AdapterInstanceID string `json:"adapterInstanceId,omitempty"`
	// A human readable status message
	Message string `json:"statusMessage,omitempty"`
	// The resource state.
	// STOPPED: Resource is not collecting
	// NOT_EXISTING: Non-existing resource
	// NONE: Resource not associated with an adapter instance.
	// UNKNOWN: Serves as a means to ensure older clients can talk to newer servers
	State string `json:"resourceState,omitempty"`
	// The resource data collection status
	// NONE: initial status of resource.
	// ERROR: error happened while collecting data
	// UNKNOWN: resource status is unknown.
	// DOWN: resource is down
	// DATA_RECEIVING: data receiving
	// OLD_DATA_RECEIVING: old data receiving
	// NO_DATA_RECEIVING: no data receiving
	// NO_PARENT_MONITORING: no parent adapter instance resource is monitoring
	// COLLECTOR_DOWN: collector is down
	Status string `json:"resourceStatus,omitempty"`
}

func unpackResourceStatusState(m interface{}) (*ResourceStatusState, error) {
	var rssm map[string]interface{}

	switch m.(type) {
	case map[string]interface{}:
		rssm = m.(map[string]interface{})
	default:
		return nil, fmt.Errorf("map not found")
	}

	r := &ResourceStatusState{}
	for k, v := range rssm {
		switch k {
		case "adapterInstanceId":
			r.AdapterInstanceID = v.(string)
		case "statusMessage":
			r.Message = v.(string)
		case "resourceState":
			r.State = v.(string)
		case "resourceStatus":
			r.Status = v.(string)
		default:
			return nil, fmt.Errorf("map contains unsupported key: %s", k)
		}
	}

	return r, nil
}
