package restapi

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// Client represents a client for interaction with a GCP REST API
type Client struct {
	Host           string
	ServiceAccount string
	Audience       string

	httpClient http.Client
}

// Do sends the API Request, parses the response as JSON, and returns the "result" value as byte
func (c *Client) Do(baseURL string, req *Request) ([]byte, error) {

	httpReq, err := req.BuildHTTPReq(c.Host, c.ServiceAccount, c.Audience, baseURL)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Print("HTTP req failed")
		return nil, err
	}

	if httpRes.StatusCode == 401 {
		return nil, errors.New("401: Unauthenticated")
	}

	defer httpRes.Body.Close()

	res, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		log.Print("HTTP decoder failed")
		return nil, err
	}

	if res == nil {
		return nil, errors.New("No result returned in REST response")
	}

	return res, nil
}
