package restapi

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// Client represents a client for interaction with a GCP REST API
type Client struct {
	Host                string
	ServiceAccount      string
	Credentials         string
	Audience            string
	Token               string
	TokenDuration       int
	TokenExpirationTime int64
	httpClient          http.Client
}

// Do sends the API Request, parses the response as JSON, and returns the HTTP status code as int, the "result" value as byte
func (c *Client) Do(baseURL string, req *Request) (int, []byte, error) {
	httpReq, err := req.BuildHTTPReq(c, baseURL)
	if err != nil {
		return 0, nil, err
	}
	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Print("HTTP req failed")
		return 0, nil, err
	}

	defer httpRes.Body.Close()

	res, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		log.Print("HTTP decoder failed")
		return 0, nil, err
	}

	if res == nil {
		return 0, nil, errors.New("No result returned in REST response")
	}

	return httpRes.StatusCode, res, nil
}
