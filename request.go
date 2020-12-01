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
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) request(method, svc string, params map[string]string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s%s%s", c.url, c.pathPrefix, svc)
	c.log.Debug(
		"making http request",
		zap.String("method", method),
		zap.String("url", reqURL),
		zap.Any("params", params),
	)

	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}

	reqURL = fmt.Sprintf("%s?%s", reqURL, q.Encode())

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

	var req *http.Request
	var err error
	req, err = http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("vRealizeOpsToken %s", c.token))
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := httpClient.Do(req)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "EOF") {
			return nil, err
		}
	}
	if res == nil {
		return nil, fmt.Errorf("response: <nil>, verify url: %s", reqURL)
	}
	defer res.Body.Close()

	c.log.Debug("http response", zap.String("status", res.Status))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("non-EOF error at url %s: %s", reqURL, err)
	}

	// c.log.Debug("http response body", zap.String("body", string(body)))

	switch res.StatusCode {
	case 200:
		return body, nil
	default:
		return nil, fmt.Errorf("error: status code %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}
