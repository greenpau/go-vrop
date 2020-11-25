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
	"strconv"
	//"go.uber.org/zap"
)

type VirtualMachineResourcesResponse struct {
	ID string `json:"_id,omitempty"`
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
	pageSize := 1

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

	/*
		if _, exists := m["attributes"]; !exists {
			return fmt.Errorf("failed to unpack VirtualMachine, attributes not found")
		}

		for k, v := range m["attributes"].(map[string]interface{}) {
			switch k {
			case "_id":
				c.ID = v.(string)
			case "device_id":
				c.DeviceID = v.(string)
			case "client_install_time":
				c.InstallTimestamp = v.(float64)
			case "client_version":
				c.Version = v.(string)
			case "users":
				for _, u := range v.([]interface{}) {
					usr := &User{}
					if err := usr.load(u.(map[string]interface{})); err != nil {
						return fmt.Errorf("failed to unpack VirtualMachine, %s attribute error: %s", k, err)
					}
					c.Users = append(c.Users, usr)
				}
			case "host_info":
				hostInfo := &Host{}
				if err := hostInfo.load(v.(map[string]interface{})); err != nil {
					return fmt.Errorf("failed to unpack VirtualMachine, %s attribute error: %s", k, err)
				}
				c.HostInfo = hostInfo
			case "last_event":
				lastEvent := &Event{}
				if err := lastEvent.load(v.(map[string]interface{})); err != nil {
					return fmt.Errorf("failed to unpack VirtualMachine, %s attribute error: %s", k, err)
				}
				c.LastEvent = lastEvent
			default:
				return fmt.Errorf("failed to unpack VirtualMachine, unsupported attribute: %s, %v", k, v)
			}
		}
	*/
	return nil
}

// UnmarshalJSON unpacks byte array into VirtualMachineResourcesResponse.
func (c *VirtualMachineResourcesResponse) UnmarshalJSON(b []byte) error {
	var requiredKeys = map[string]bool{
		"resourceList": false,
		"pageInfo":     false,
	}
	var m map[string]interface{}
	if len(b) < 10 {
		return fmt.Errorf("invalid VirtualMachineResourcesResponse data: %s", b)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unpack VirtualMachineResourcesResponse")
	}

	return fmt.Errorf("YYY: %s", string(b))
}
