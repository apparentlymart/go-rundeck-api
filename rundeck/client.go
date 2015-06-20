// Package rundeck provides a client for interacting with a Rundeck instance
// via its HTTP API.
//
// Instantiate a Client with the NewClient function to get started.
//
// At present this package uses Rundeck API version 13.
package rundeck

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"mime/multipart"
	"strings"
)

// ClientConfig is used with NewClient to specify initialization settings.
type ClientConfig struct {
	// The base URL of the Rundeck instance.
	BaseURL string

	// The API auth token generated from user settings in the Rundeck UI.
	AuthToken string

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

	apiPath, _ := url.Parse("api/13/")
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
		req.ContentLength = int64(len(reqBodyBytes))
		req.Header.Add("Content-Type", "application/xml")
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
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

	if result != nil {
		if res.StatusCode != 200 && res.StatusCode != 201 {
			return fmt.Errorf("Server did not return an XML payload")
		}
		err = xml.Unmarshal(resBodyBytes, result)
		if err != nil {
			return fmt.Errorf("Error decoding response XML payload: %s", err.Error())
		}
	}

	return nil
}

func (c *Client) get(pathParts []string, query map[string]string, result interface{}) error {
	return c.request("GET", pathParts, query, nil, result)
}

func (c *Client) post(pathParts []string, query map[string]string, reqBody interface{}, result interface{}) error {
	return c.request("POST", pathParts, query, reqBody, result)
}

func (c *Client) put(pathParts []string, reqBody interface{}, result interface{}) error {
	return c.request("PUT", pathParts, nil, reqBody, result)
}

func (c *Client) delete(pathParts []string) error {
	return c.request("DELETE", pathParts, nil, nil, nil)
}

func (c *Client) postXMLBatch(pathParts []string, args map[string]string, xmlBatch interface{}, result interface{}) error {
	req := &http.Request{
		Method: "POST",
		Header: http.Header{},
	}
	req.Header.Add("User-Agent", "Go-Rundeck-API")
	req.Header.Add("X-Rundeck-Auth-Token", c.authToken)

	urlPath := &url.URL{
		Path: strings.Join(pathParts, "/"),
	}
	reqURL := c.apiURL.ResolveReference(urlPath)
	req.URL = reqURL

	buf := bytes.Buffer{}
	writer := multipart.NewWriter(&buf)
	for k, v := range args {
		err := writer.WriteField(k, v)
		if err != nil {
			return err
		}
	}
	partWriter, err := writer.CreateFormFile("xmlBatch", "batch.xml")
	if err != nil {
		return err
	}

	reqBodyBytes, err := xml.Marshal(xmlBatch)
	if err != nil {
		return err
	}

	_, err = partWriter.Write(reqBodyBytes)
	if err != nil {
		return err
	}

	writer.Close()

	reqBodyReader := bytes.NewReader(buf.Bytes())
	req.Body = ioutil.NopCloser(reqBodyReader)
	req.ContentLength = int64(buf.Len())
	req.Header.Add("Content-Type", writer.FormDataContentType())

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
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

	if result != nil {
		if res.StatusCode != 200 && res.StatusCode != 201 {
			return fmt.Errorf("Server did not return an XML payload")
		}
		err = xml.Unmarshal(resBodyBytes, result)
		if err != nil {
			return fmt.Errorf("Error decoding response XML payload: %s", err.Error())
		}
	}

	return nil
}
