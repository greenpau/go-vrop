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
	"go.uber.org/zap"
	"strconv"
)

type VirtualMachineResourcesResponse struct {
	Page      *PageInfo   `json:"pageInfo,omitempty"`
	Links     []*Link     `json:"links,omitempty"`
	Resources []*Resource `json:"resourceList,omitempty"`
}

// VirtualMachine is a virtual machine.
type VirtualMachine struct {
	ID string `json:"_id,omitempty"`
}

// GetVirtualMachines returns a list of VirtualMachine instances.
func (c *Client) GetVirtualMachines(opts map[string]interface{}) ([]*VirtualMachine, error) {
	machines := []*VirtualMachine{}
	if err := c.authenticate(); err != nil {
		return machines, err
	}

	pageOffset := 0
	pageSize := 100

	params := make(map[string]string)
	params["resourceKind"] = "virtualmachine"
	params["page"] = strconv.Itoa(pageOffset)
	params["pageSize"] = strconv.Itoa(pageSize)
	b, err := c.request("GET", "resources", params)
	if err != nil {
		return machines, err
	}

	//return machines, fmt.Errorf(string(b))

	resp := &VirtualMachineResourcesResponse{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return machines, fmt.Errorf("failed unmarshalling response: %s", err)
	}

	c.log.Warn("XXXX", zap.Any("XXXX", resp))

	//c.log.Warn("response", zap.Any("response", string(b)))
	//if resp.Status != "success" {
	//	return machines, fmt.Errorf("failed request: %s", resp.Message)
	//}
	//for _, machine := range resp.VirtualMachines {
	//c.log.Warn("machine", zap.Any("machine", machine))
	//	machines = append(machines, machine)
	//}

	return machines, nil
}

// ToJSONString serializes VirtualMachine to a string.
func (c *VirtualMachine) ToJSONString() (string, error) {
	itemJSON, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed converting to json: %s", err)
	}
	return string(itemJSON), nil
}

// UnmarshalJSON unpacks byte array into VirtualMachine.
func (c *VirtualMachine) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if len(b) < 10 {
		return fmt.Errorf("invalid VirtualMachine data: %s", b)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unpack VirtualMachine")
	}

	return fmt.Errorf("XXX: %s", string(b))

	// return nil
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

	if p, err := unpackPageInfo(m["pageInfo"]); err != nil {
		return fmt.Errorf("failed to unpack %s pageInfo: %s", obj, err)
	} else {
		c.Page = p
	}

	for _, item := range m["links"].([]interface{}) {
		if link, err := unpackLink(item); err != nil {
			return fmt.Errorf("failed to unpack %s link: %s", obj, err)
		} else {
			c.Links = append(c.Links, link)
		}
	}

	for _, item := range m["resourceList"].([]interface{}) {
		if resource, err := unpackResource(item); err != nil {
			return fmt.Errorf("failed to unpack %s resourceList: %s", obj, err)
		} else {
			c.Resources = append(c.Resources, resource)
		}
	}

	//return fmt.Errorf("YYY: %s", string(b))
	//return fmt.Errorf("YYY: %v", c.Resources[0])
	return nil
}
