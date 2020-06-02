package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"log"
)

// operateActiveDirectoryRequest requests the user's input for creating/updating an active directory
type operateActiveDirectoryRequest struct {
	Username string `structs:"username"`
	Password string `structs:"password"`
	Region   string `structs:"region"`
	Domain   string `structs:"domain"`
	DNS      string `structs:"DNS"`
	NetBIOS  string `structs:"netBIOS"`
	UUID     string `structs:"UUID"`
}

// operateActiveDirectoryResult returns the api response for creating/updating an active directory
type operateActiveDirectoryResult struct {
	UUID   string `json:"UUID"`
	Region string `json:"region"`
}

// listActiveDirectoryRequest requests the region and uuid of the active directory being fetched
type listActiveDirectoryRequest struct {
	Region string `structs:"region"`
	UUID   string `structs:"UUID"`
}

// listActiveDirectoryResult lists the active directory for given ID
type listActiveDirectoryResult struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Region   string `json:"region"`
	Domain   string `json:"domain"`
	DNS      string `json:"DNS"`
	NetBIOS  string `json:"netBIOS"`
	UUID     string `json:"UUID"`
}

type listActiveDirectoryApiResult struct {
	Collection []listActiveDirectoryResult
}

// deleteActiveDirectoryRequest requests the region and uuid of the active directory being deleted
type deleteActiveDirectoryRequest struct {
	Region string `structs:"region"`
	UUID   string `structs:"UUID"`
}

func (c *Client) createActiveDirectory(request *operateActiveDirectoryRequest) (operateActiveDirectoryResult, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Storage/ActiveDirectory", request.Region)
	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("CreateActiveDirectory request failed")
		return operateActiveDirectoryResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "CreateActiveDirectory")
	if responseError != nil {
		return operateActiveDirectoryResult{}, responseError
	}

	var result operateActiveDirectoryResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from CreateActiveDirectory")
		return operateActiveDirectoryResult{}, err
	}

	return result, nil
}

func (c *Client) listActiveDirectoryForRegion(request listActiveDirectoryRequest) (listActiveDirectoryResult, error) {
	// GCP only allows one active directory per region.
	baseURL := fmt.Sprintf("%s/Storage/ActiveDirectory", request.Region)
	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("listActiveDirectory request failed")
		return listActiveDirectoryResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "listActiveDirectory")
	if responseError != nil {
		return listActiveDirectoryResult{}, responseError
	}

	var active_directorys []listActiveDirectoryResult
	if err := json.Unmarshal(response, &active_directorys); err != nil {
		log.Print("Failed to unmarshall response from listActiveDirectory")
		return listActiveDirectoryResult{}, err
	}
	for _, v := range active_directorys {
		// only one active directory is allowed in each region. Region is the unique identifier if uuid doesn't exist yet.
		if v.Region == request.Region {
			return v, nil
		}

	}

	return listActiveDirectoryResult{}, nil
}

func (c *Client) deleteActiveDirectory(request deleteActiveDirectoryRequest) error {
	baseURL := fmt.Sprintf("%s/Storage/ActiveDirectory/%s", request.Region, request.UUID)
	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("deleteActiveDirectory request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteActiveDirectory")
	if responseError != nil {
		return responseError
	}

	return nil
}

func (c *Client) updateActiveDirectory(request operateActiveDirectoryRequest) error {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Storage/ActiveDirectory/%s", request.Region, request.UUID)
	statusCode, response, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Print("updateActiveDirectory request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "updateActiveDirectory")
	if responseError != nil {
		return responseError
	}

	var result listActiveDirectoryResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from updateActiveDirectory")
		return err
	}

	return nil
}
