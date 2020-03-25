package gcp

import "fmt"

// Config is a struct for user input
type configStuct struct {
	Project        string
	ServiceAccount string
}

// Client is the main function to connect to the APi
func (c *configStuct) clientFun() (*Client, error) {
	client := &Client{
		Host:     fmt.Sprintf("https://cloudvolumesgcp-api.netapp.com/v2/projects/%s/locations/", c.Project),
		Audience: "https://cloudvolumesgcp-api.netapp.com",
	}

	client.SetServiceAccount(c.ServiceAccount)
	client.SetProjectID(c.Project)

	return client, nil
}
