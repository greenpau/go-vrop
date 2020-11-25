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
	"fmt"
	"go.uber.org/zap"
	"time"
)

// ReceiverDataLimit is the limit of data in bytes the client will read
// from a server.
const ReceiverDataLimit int64 = 1e6

// Client is an instance of Proofpoint API client.
type Client struct {
	url                string
	host               string
	port               int
	protocol           string
	username           string
	password           string
	token              string
	tokenExpiresAt     time.Time
	validateServerCert bool
	dataLimit          int64
	pathPrefix         string
	log                *zap.Logger
}

// NewClient returns an instance of Client.
func NewClient(opts map[string]interface{}) (*Client, error) {
	c := &Client{
		host:       "vrop",
		port:       443,
		protocol:   "https",
		pathPrefix: "/suite-api/api/",
		dataLimit:  ReceiverDataLimit,
	}
	log, err := newLogger(opts)
	if err != nil {
		return nil, fmt.Errorf("failed initializing log: %s", err)
	}
	c.log = log
	return c, nil
}

// Close performs a cleanup associated with Client..
func (c *Client) Close() {
	if c.log != nil {
		c.log.Sync()
	}
}

// Info sends information about Client to the configured logger.
func (c *Client) Info() {
	c.rebaseURL()
	c.log.Debug(
		"client configuration",
		zap.String("url", c.url),
		zap.String("path_prefix", c.pathPrefix),
	)
}

func (c *Client) rebaseURL() {
	if (c.protocol == "https" && c.port == 443) ||
		(c.protocol == "http" && c.port == 80) {
		c.url = fmt.Sprintf("%s://%s", c.protocol, c.host)
		return
	}
	c.url = fmt.Sprintf("%s://%s:%d", c.protocol, c.host, c.port)
	return
}

// SetHost sets the target host for the API calls.
func (c *Client) SetHost(s string) error {
	if s == "" {
		return fmt.Errorf("empty hostname or ip address")
	}
	c.host = s
	c.rebaseURL()
	return nil
}

// SetPort sets the port number for the API calls.
func (c *Client) SetPort(p int) error {
	if p == 0 {
		return fmt.Errorf("invalid port: %d", p)
	}
	c.port = p
	c.rebaseURL()
	return nil
}

// SetUsername sets API username
func (c *Client) SetUsername(s string) error {
	if s == "" {
		return fmt.Errorf("empty username")
	}
	c.username = s
	return nil
}

// SetPassword sets API password.
func (c *Client) SetPassword(s string) error {
	if s == "" {
		return fmt.Errorf("empty password")
	}
	c.password = s
	return nil
}

// SetProtocol sets the protocol for the API calls.
func (c *Client) SetProtocol(s string) error {
	switch s {
	case "http":
		c.protocol = s
	case "https":
		c.protocol = s
	default:
		return fmt.Errorf("supported protocols: http, https; unsupported protocol: %s", s)
	}
	c.rebaseURL()
	return nil
}

// SetValidateServerCertificate instructs the client to enforce the validation of certificates
// and check certificate errors.
func (c *Client) SetValidateServerCertificate() error {
	c.validateServerCert = true
	return nil
}
