package gcp

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
)

// operateActiveDirectoryRequest requests the user's input for creating/updating an active directory
type operateActiveDirectoryRequest struct {
	Username                   string   `structs:"username"`
	Password                   string   `structs:"password"`
	Region                     string   `structs:"region"`
	Domain                     string   `structs:"domain"`
	DNS                        string   `structs:"DNS"`
	NetBIOS                    string   `structs:"netBIOS"`
	OrganizationalUnit         string   `structs:"organizationalUnit"`
	Site                       string   `structs:"site"`
	UUID                       string   `structs:"UUID"`
	LdapSigning                bool     `structs:"ldapSigning"`
	KdcIP                      string   `structs:"kdcIP"`
	AllowLocalNFSUsersWithLdap bool     `structs:"allowLocalNFSUsersWithLdap"`
	SecurityOperators          []string `structs:"securityOperators"`
	BackupOperators            []string `structs:"backupOperators"`
	AesEncryption              bool     `structs:"aesEncryption"`
	Label                      string   `structs:"label"`
	AdName                     string   `structs:"adName"`
	ManagedAD                  bool     `structs:"managedAD"`
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
	Username                   string   `json:"username"`
	Password                   string   `json:"password"`
	Region                     string   `json:"region"`
	Domain                     string   `json:"domain"`
	DNS                        string   `json:"DNS"`
	NetBIOS                    string   `json:"netBIOS"`
	OrganizationalUnit         string   `structs:"organizationalUnit"`
	Site                       string   `structs:"site"`
	UUID                       string   `json:"UUID"`
	LdapSigning                bool     `json:"ldapSigning"`
	KdcIP                      string   `json:"kdcIP"`
	AllowLocalNFSUsersWithLdap bool     `json:"allowLocalNFSUsersWithLdap"`
	SecurityOperators          []string `json:"securityOperators"`
	BackupOperators            []string `json:"backupOperators"`
	AesEncryption              bool     `json:"aesEncryption"`
	Label                      string   `json:"label"`
	AdName                     string   `json:"adName"`
	ManagedAD                  bool     `structs:"managedAD"`
}

type listActiveDirectoryAPIResult struct {
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

	//code for terraform import
	var originalID string = ""
	s := strings.Split(request.UUID, ":")
	if len(s) == 2 {
		originalID = request.UUID
		request.UUID = s[0]
		request.Region = s[1]
	}

	if request.Region == "" {
		// terraform import: ID = <activeDirectoryID> and no region specified
		// find all activeDirectories which match activeDirectoryID
		activeDirectory, err := c.filterAllActiveDirectories(func(v listActiveDirectoryResult) bool {
			return v.UUID == request.UUID
		})
		if err != nil {
			return listActiveDirectoryResult{}, err
		}

		if len(activeDirectory) == 0 {
			return listActiveDirectoryResult{}, fmt.Errorf("listActiveDirectoryForRegion: No active directory found with ID %s", request.UUID)
		}
		if len(activeDirectory) > 1 {
			// return error message which tells user to rerun with ID = <volumeID>:<region> format
			return listActiveDirectoryResult{}, fmt.Errorf(`listActiveDirectoryForRegion: More than one active directory found with ID %s. \n
			If this happend while running terraform import, please use \n
			terraform import ADDR ID, with ID using <active directory>:<region_name> format`, request.UUID)
		}
		if len(activeDirectory) == 1 {
			// we found the right activeDirectory to import
			request.Region = activeDirectory[0].Region
		}
	}

	baseURL := fmt.Sprintf("%s/Storage/ActiveDirectory", request.Region)
	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("listActiveDirectory request failed")
		return listActiveDirectoryResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "listActiveDirectoryForRegion")
	if responseError != nil {
		return listActiveDirectoryResult{}, responseError
	}

	var activeDirectories []listActiveDirectoryResult
	if err := json.Unmarshal(response, &activeDirectories); err != nil {
		log.Print("Failed to unmarshall response from listActiveDirectoryForRegion")
		return listActiveDirectoryResult{}, err
	}

	for _, v := range activeDirectories {
		// only one active directory is allowed in each region. Region is the unique identifier if uuid doesn't exist yet.
		if v.Region == request.Region {
			if originalID != "" {
				v.UUID = originalID
			}
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

// Returns volumes of the project. region = "-" for all regions
func (c *Client) getActiveDirectories(region string) ([]listActiveDirectoryResult, error) {

	baseURL := fmt.Sprintf("%s/Storage/ActiveDirectory", region)
	var result []listActiveDirectoryResult

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("getActiveDirectories request failed")
		return result, err
	}

	responseError := apiResponseChecker(statusCode, response, "getActiveDirectories")
	if responseError != nil {
		return result, responseError
	}

	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getActiveDirectories")
		return result, err
	}
	return result, nil
}

// Filter all volumes of the project by applying a filter function
// Example filter function: func(v volumeResult) bool { return v.VolumeID == "1234-5678-90" }
func (c *Client) filterAllActiveDirectories(f func(listActiveDirectoryResult) bool) ([]listActiveDirectoryResult, error) {
	filteredActiveDirectories := make([]listActiveDirectoryResult, 0)

	vols, err := c.getActiveDirectories("-")
	if err != nil {
		return filteredActiveDirectories, err
	}

	for _, v := range vols {
		if f(v) {
			filteredActiveDirectories = append(filteredActiveDirectories, v)
		}
	}
	return filteredActiveDirectories, nil
}
