package gcp

import (
	"sync"

	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
	"github.com/sirupsen/logrus"
)

var ourlog = logrus.WithFields(logrus.Fields{
	"prefix": "main",
})

// A Client to interact with the GCP REST API
type Client struct {
	Host                  string
	MaxConcurrentRequests int
	BaseURL               string
	ServiceAccount        string
	Project               string
	Audience              string

	initOnce      sync.Once
	restapiClient *restapi.Client
	requestSlots  chan int
}

// CallAPIMethod can be used to make a request to any GCP API method, receiving results as byte
func (c *Client) CallAPIMethod(method string, baseURL string, params map[string]interface{}) (int, []byte, error) {
	c.initOnce.Do(c.init)

	c.waitForAvailableSlot()
	defer c.releaseSlot()

	ourlog.WithFields(logrus.Fields{
		"method": method,
		"params": params,
	}).Debug("Calling API")

	if params == nil {
		params = map[string]interface{}{}
	}
	statusCode, result, err := c.restapiClient.Do(baseURL, &restapi.Request{
		Method: method,
		Params: params,
	})
	if err != nil {
		return statusCode, nil, err
	}
	ourlog.WithFields(logrus.Fields{
		"method": method,
	}).Debug("Received successful API response")
	return statusCode, result, nil
}

func (c *Client) init() {
	if c.MaxConcurrentRequests == 0 {
		c.MaxConcurrentRequests = 6
	}
	c.requestSlots = make(chan int, c.MaxConcurrentRequests)
	c.restapiClient = &restapi.Client{
		Host:           c.Host,
		ServiceAccount: c.ServiceAccount,
		Audience:       c.Audience,
	}
}

// SetServiceAccount for the client to use for requests to the GCP API
func (c *Client) SetServiceAccount(serviceAccount string) {
	c.ServiceAccount = serviceAccount
}

// GetServiceAccount returns the API version that will be used for GCP API requests
func (c *Client) GetServiceAccount() string {
	return c.ServiceAccount
}

func (c *Client) waitForAvailableSlot() {
	c.requestSlots <- 1
}

func (c *Client) releaseSlot() {
	<-c.requestSlots
}
