// Package rundeck provides a client for interacting with a Rundeck instance
// via its HTTP API.
//
// Instantiate a Client with the NewClient function to get started.
//
// At present this package uses Rundeck API version 12.
package rundeck

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"bytes"
)

// ClientConfig is used with NewClient to specify initialization settings.
type ClientConfig struct {
	// The base URL of the Rundeck instance.
	BaseURL            string

	// The API auth token generated from user settings in the Rundeck UI.
	AuthToken          string

	// Don't fail if the server uses SSL with an un-verifiable certificate.
	// This is not recommended except during development/debugging.
	AllowUnverifiedSSL bool
}

// Client is a Rundeck API client interface.
type Client struct {
	httpClient *http.Client
	apiURL     *url.URL
	authToken  string
}


// NewClient returns a configured Rundeck client.
func NewClient(config *ClientConfig) (*Client, error) {
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.AllowUnverifiedSSL,
		},
	}
	httpClient := &http.Client{
		Transport: t,
	}

	apiPath, _ := url.Parse("api/12/")
	baseUrl, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("Invalid base URL: %s", err.Error())
	}
	apiURL := baseUrl.ResolveReference(apiPath)

	return &Client{
		httpClient: httpClient,
		apiURL:     apiURL,
		authToken:  config.AuthToken,
	}, nil
}

func (c *Client) request(method string, pathParts []string, query map[string]string, reqBody interface{}, result interface{}) error {

	req := &http.Request{
		Method: method,
		Header: http.Header{},
	}
	req.Header.Add("User-Agent", "Go-Rundeck-API")
	req.Header.Add("X-Rundeck-Auth-Token", c.authToken)
	req.Header.Add("Accept", "application/xml")

	urlPath := &url.URL{
		Path: strings.Join(pathParts, "/"),
	}
	reqURL := c.apiURL.ResolveReference(urlPath)
	req.URL = reqURL

	if len(query) > 0 {
		urlQuery := url.Values{}
		for k, v := range query {
			urlQuery.Add(k, v)
		}
		reqURL.RawQuery = urlQuery.Encode()
	}

	if reqBody != nil {
		reqBodyBytes, err := xml.Marshal(reqBody)
		if err != nil {
			return err
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		if strings.HasPrefix(res.Header.Get("Content-Type"), "text/xml") {
			var richErr Error
			err = xml.Unmarshal(resBodyBytes, &richErr)
			if err != nil {
				return fmt.Errorf("HTTP Error %i with error decoding XML body: %s", res.StatusCode, err.Error())
			}
			return richErr
		} else {
			return fmt.Errorf("HTTP Error %i", res.StatusCode)
		}
	}

	err = xml.Unmarshal(resBodyBytes, result)
	if err != nil {
		err = fmt.Errorf("Error decoding response XML payload: %s", err.Error())
	}

	return err
}

func (c *Client) get(pathParts []string, query map[string]string, result interface{}) error {
	return c.request("GET", pathParts, query, nil, result)
}
