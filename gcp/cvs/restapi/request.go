package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
)

// Request represents a request to a REST API
type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// BuildHTTPReq builds an HTTP request to carry out the REST request
func (r *Request) BuildHTTPReq(host string, serviceAccount string, audience string, baseURL string) (*http.Request, error) {
	bodyJSON, err := json.Marshal(r.Params)
	if err != nil {
		return nil, err
	}

	url := host + baseURL
	req, err := http.NewRequest(r.Method, url, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}

	keyBytes, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("Unable to read service account key file  %v", err)
	}
	tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, audience)
	if err != nil {
		return nil, fmt.Errorf("Error building JWT access token source: %v", err)
	}
	jwt, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("Unable to generate JWT token: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt.AccessToken)

	return req, nil
}
