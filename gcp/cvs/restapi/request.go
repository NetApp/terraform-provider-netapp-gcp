package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"golang.org/x/oauth2/google"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

// Request represents a request to a REST API
type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// BuildHTTPReq builds an HTTP request to carry out the REST request
func (r *Request) BuildHTTPReq(c *Client, baseURL string) (*http.Request, error) {
	var keyBytes []byte
	var err error
	var req *http.Request
	url := c.Host + baseURL
	if r.Method != "GET" && r.Method != "DELETE" {
		bodyJSON, err := json.Marshal(r.Params)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(r.Method, url, bytes.NewReader(bodyJSON))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(r.Method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	// Can be specified in multiple ways:
	// 1. JSON key as base64-encoded string - credentials
	// 2. Service Account principal name when using service account impersonation - service_account
	// 3. Absolute file path to an JSON key file - service_account
	token := c.Token
	expirationTime := c.TokenExpirationTime
	re := regexp.MustCompile(`^[[a-z]([-a-z0-9]*[a-z0-9])@[a-z0-9-]+\.iam\.gserviceaccount\.com$`)
	if c.Credentials != "" {
		keyBytes = []byte(c.Credentials)
	} else if c.ServiceAccount != "" {
		if re.MatchString(c.ServiceAccount) {
			// Use existing token, unless it is expired
			if time.Now().Unix() >= c.TokenExpirationTime {
				token, expirationTime, err = getToken(c.ServiceAccount, c.TokenDuration)
				if err != nil {
					return nil, fmt.Errorf("Unable to get token from %s %v", c.ServiceAccount, err)
				}
				log.Printf("Update token %v and expiration time %v", token, expirationTime)
				c.TokenExpirationTime = expirationTime
				c.Token = token
			}
		} else {
			keyBytes, err = ioutil.ReadFile(c.ServiceAccount)
			if err != nil {
				return nil, fmt.Errorf("Unable to read service account key file  %v", err)
			}
		}
	} else {
		return nil, fmt.Errorf("Need credential or service_account to get the authentication")
	}

	if keyBytes != nil {
		tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, c.Audience)
		if err != nil {
			return nil, fmt.Errorf("Error building JWT access token source: %v", err)
		}
		jwt, err := tokenSource.Token()
		if err != nil {
			return nil, fmt.Errorf("Unable to generate JWT token: %v", err)
		}
		token = jwt.AccessToken
		c.Token = token
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

func getToken(serviceAccountName string, tokenDuration int) (string, int64, error) {
	log.Printf("getToken...")
	if tokenDuration <= 0 || tokenDuration > 60 {
		log.Print("tokenDuration is set to 60 min")
		tokenDuration = 60
	}
	ctx := context.Background()
	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("Get iam client err: %v", err)
	}
	defer c.Close()

	expTime := time.Now().Add(time.Minute * time.Duration(tokenDuration)).Unix()
	jwtPayload := map[string]interface{}{
		"iss": serviceAccountName,
		"iat": time.Now().Unix(),
		"aud": "https://cloudvolumesgcp-api.netapp.com/",
		"sub": serviceAccountName,
		"exp": expTime,
	}
	log.Printf("jwtPayload: %v\n", jwtPayload)
	payloadBytes, err := json.Marshal(jwtPayload)
	req := &credentialspb.SignJwtRequest{
		// See https://pkg.go.dev/google.golang.org/genproto/googleapis/iam/credentials/v1#SignJwtRequest.
		Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountName),
		Payload: string(payloadBytes),
	}
	resp, err := c.SignJwt(ctx, req)
	if err != nil {
		return "", 0, fmt.Errorf("signjwt failed: %v", err)
	}

	log.Printf("SignedJwt: %v\n", resp.SignedJwt)
	return resp.SignedJwt, expTime, nil
}
