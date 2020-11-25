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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// AuthResponse is the response from auth token acquisition endpoint.
type AuthResponse struct {
	Token     string        `json:"token,omitempty"`
	Validity  float64       `json:"validity,omitempty"`
	ExpiresAt string        `json:"expiresAt,omitempty"`
	Roles     []interface{} `json:"roles,omitempty"`
}

func (c *Client) authenticate() error {
	if c.token != "" {
		return nil
	}
	reqURL := fmt.Sprintf("%s%sauth/token/acquire", c.url, c.pathPrefix)
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	if !c.validateServerCert {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}

	data := map[string]string{
		"username": c.username,
		"password": c.secret,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("authentication error: %s", err)
	}

	c.log.Debug("http request", zap.String("url", reqURL))

	var req *http.Request
	req, err = http.NewRequest("POST", reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := httpClient.Do(req)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "EOF") {
			return err
		}
	}
	if res == nil {
		return fmt.Errorf("response: <nil>, verify url: %s", reqURL)
	}
	defer res.Body.Close()

	c.log.Debug("http response", zap.String("status", res.Status))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("non-EOF error at url %s: %s", reqURL, err)
	}

	c.log.Debug("http response body", zap.String("body", string(body)))

	if res.StatusCode != 200 {
		return fmt.Errorf("authentication failed")
	}

	authResp := &AuthResponse{}
	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("failed unmarshalling authentication response: %s", err)
	}

	if authResp.Token == "" {
		return fmt.Errorf("token not found in authentication response")
	}

	c.token = authResp.Token

	secs := int64(authResp.Validity / 1000)
	nsecs := int64(((authResp.Validity / 1000) - float64(secs)) * 1e9)
	c.tokenExpiresAt = time.Unix(secs, nsecs)

	c.log.Debug(
		"authenticated successfully",
		zap.String("token_expires_at", c.tokenExpiresAt.String()),
	)

	return nil
}
