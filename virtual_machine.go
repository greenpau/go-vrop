// Copyright 2020 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vrop

import (
	"encoding/json"
	"fmt"
	//"go.uber.org/zap"
	"strconv"
	"time"
)

// VirtualMachineResourcesResponse is a response with VirtualMachine resources.
type VirtualMachineResourcesResponse struct {
	Page      *PageInfo   `json:"pageInfo,omitempty"`
	Links     []*Link     `json:"links,omitempty"`
	Resources []*Resource `json:"resourceList,omitempty"`
}

// VirtualMachine is a virtual machine.
type VirtualMachine struct {
	ID                         string    `json:"id,omitempty"`
	Name                       string    `json:"name,omitempty"`
	VMEntityInstanceUUID       string    `json:vm_entity_instance_uuid,omitempty"`
	VMEntityName               string    `json:vm_entity_name,omitempty"`
	VMEntityObjectID           string    `json:vm_entity_object_id,omitempty"`
	VMEntityVCID               string    `json:vm_entity_vcid,omitempty"`
	VMServiceMonitoringEnabled bool      `json:vm_service_monitoring_enabled,omitempty"`
	CreatedAt                  time.Time `json:created_at,omitempty"`
	LastSeenAt                 time.Time `json:last_seen_at,omitempty"`
	Errors                     []error   `json:errors,omitempty"`
}

// GetVirtualMachines returns a list of VirtualMachine instances.
func (c *Client) GetVirtualMachines(opts map[string]interface{}) ([]*VirtualMachine, error) {
	machines := []*VirtualMachine{}
	if err := c.authenticate(); err != nil {
		return machines, err
	}

	pageOffset := 0
	pageSize := 10

	for {
		params := make(map[string]string)
		params["resourceKind"] = "virtualmachine"
		params["page"] = strconv.Itoa(pageOffset)
		params["pageSize"] = strconv.Itoa(pageSize)
		b, err := c.request("GET", "resources", params)
		if err != nil {
			return machines, err
		}

		// return machines, fmt.Errorf(string(b))

		resp := &VirtualMachineResourcesResponse{}
		if err := json.Unmarshal(b, &resp); err != nil {
			return machines, fmt.Errorf("failed unmarshalling response: %s", err)
		}

		for _, r := range resp.Resources {
			m := &VirtualMachine{}
			m.ID = r.ID
			m.CreatedAt = r.CreationTime
			m.LastSeenAt = time.Now().UTC()
			m.Name = r.Key.Name
			for _, entry := range r.Key.ResourceIdentifiers {
				switch entry.Key {
				case "VMEntityInstanceUUID":
					m.VMEntityInstanceUUID = entry.Value
				case "VMEntityName":
					m.VMEntityName = entry.Value
				case "VMEntityObjectID":
					m.VMEntityObjectID = entry.Value
				case "VMEntityVCID":
					m.VMEntityVCID = entry.Value
				case "VMServiceMonitoringEnabled":
					if entry.Value == "true" || entry.Value == "True" || entry.Value == "TRUE" {
						m.VMServiceMonitoringEnabled = true
					}
				}
			}
			if err := m.GetProperties(); err != nil {
				m.Errors = append(m.Errors, err)

			}
			machines = append(machines, m)
		}

		if len(resp.Resources) < pageSize {
			break
		}
		pageOffset++
		// TODO: remove break
		break
	}

	return machines, nil
}

// ToJSONString serializes VirtualMachine to a string.
func (m *VirtualMachine) ToJSONString() (string, error) {
	itemJSON, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed converting to json: %s", err)
	}
	return string(itemJSON), nil
}

// UnmarshalJSON unpacks byte array into VirtualMachineResourcesResponse.
func (c *VirtualMachineResourcesResponse) UnmarshalJSON(b []byte) error {
	obj := "VirtualMachineResourcesResponse"
	var requiredKeys = map[string]bool{
		"resourceList": false,
		"pageInfo":     false,
		"links":        false,
	}
	var optionalKeys = map[string]bool{}
	var m map[string]interface{}
	if len(b) < 10 {
		return fmt.Errorf("invalid %s data: %s", obj, b)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unpack %s", obj)
	}

	for k := range m {
		if _, exists := requiredKeys[k]; exists {
			requiredKeys[k] = true
			continue
		}
		if _, exists := optionalKeys[k]; exists {
			optionalKeys[k] = true
			continue
		}
		return fmt.Errorf("failed to unpack %s, found unsupported key: %s", obj, k)
	}

	for k, present := range requiredKeys {
		if !present {
			return fmt.Errorf("failed to unpack %s, required key not found: %s", obj, k)
		}
	}

	p, err := unpackPageInfo(m["pageInfo"])
	if err != nil {
		return fmt.Errorf("failed to unpack %s pageInfo: %s", obj, err)
	}
	c.Page = p

	for _, item := range m["links"].([]interface{}) {
		link, err := unpackLink(item)
		if err != nil {
			return fmt.Errorf("failed to unpack %s link: %s", obj, err)
		}
		c.Links = append(c.Links, link)
	}

	for _, item := range m["resourceList"].([]interface{}) {
		resource, err := unpackResource(item)
		if err != nil {
			return fmt.Errorf("failed to unpack %s resourceList: %s", obj, err)
		}
		c.Resources = append(c.Resources, resource)
	}

	return nil
}

// GetProperties fetches latest properties of VirtualMachine.
func (m *VirtualMachine) GetProperties() error {
	return nil
}
