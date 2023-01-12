package gcp

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

type storagePool struct {
	Name                   string         `json:"name"`
	Network                string         `json:"network"`
	Region                 string         `json:"region"`
	ServiceLevel           string         `json:"serviceLevel"`
	SizeInBytes            int            `json:"sizeInBytes"`
	RegionalHA             bool           `json:"regionalHA"`
	GlobalILB              bool           `json:"globalILB"`
	ManagedPool            bool           `json:"managedPool"`
	SecondaryZone          string         `json:"secondaryZone"`
	Zone                   string         `json:"zone"`
	PoolID                 string         `json:"poolId"`
	StorageClass           string         `json:"storageClass"`
	Jobs                   []job          `json:"jobs"`
	BillingLabels          []billingLabel `json:"billingLabels"`
	State                  string         `json:"state"`
	SharedVpcProjectNumber string
}

func (c *Client) createStoragePool(request *storagePool) (storagePool, error) {
	var projectID string
	if request.SharedVpcProjectNumber != "" {
		projectID = request.SharedVpcProjectNumber
	} else {
		projectID = c.GetProjectID()
	}
	request.Network = fmt.Sprintf("projects/%s/global/networks/%s", projectID, request.Network)
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Pools", request.Region)
	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Printf("createStoragePool request failed: %#v", err)
		return storagePool{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "createStoragePool")
	if responseError != nil {
		return storagePool{}, responseError
	}
	var contentHolder map[string]interface{}
	if err := json.Unmarshal(response, &contentHolder); err != nil {
		log.Printf("Failed to unmarshall response from createStoragePool: %#v", err)
		return storagePool{}, err
	}
	responseHolder := contentHolder["response"].(map[string]interface{})
	anyValueHolder := responseHolder["AnyValue"].(map[string]interface{})
	poolData, err := json.Marshal(anyValueHolder)
	if err != nil {
		return storagePool{}, err
	}
	var result storagePool
	if err := json.Unmarshal(poolData, &result); err != nil {
		log.Printf("Failed to unmarshall response from createStoragePool: %#v", err)
		return storagePool{}, err
	}
	err = c.waitForJobCompletion(result.Region, result.Jobs[0].JobID, 1200, 20, false)
	if err != nil {
		return storagePool{}, err
	}
	return result, nil
}

func (c *Client) getStoragePools(request *storagePool) ([]storagePool, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Pools", request.Region)
	var result []storagePool
	statusCode, response, err := c.CallAPIMethod("GET", baseURL, params)
	if err != nil {
		log.Printf("getStoragePools request failed: %#v", err)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getStoragePools")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Printf("Failed to unmarshall response from getStoragePools: %#v", err)
		return result, err
	}
	return result, nil
}

func (c *Client) getStoragePoolByID(request *storagePool) (storagePool, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Pools/%s", request.Region, request.PoolID)
	var result storagePool
	statusCode, response, err := c.CallAPIMethod("GET", baseURL, params)
	if err != nil {
		log.Printf("getStoragePoolByID request failed: %#v", err)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getStoragePoolByID")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Printf("Failed to unmarshall response from getStoragePoolByID: %#v", err)
		return result, err
	}
	return result, nil
}

func (c *Client) deleteStoragePool(request *storagePool) error {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Pools/%s", request.Region, request.PoolID)
	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, params)
	if err != nil {
		log.Printf("deleteStoragePool request failed: %#v", err)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteStoragePool")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) updateStoragePool(request *storagePool) error {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Pools/%s", request.Region, request.PoolID)
	statusCode, response, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Printf("updateStoragePool request failed: %#v", err)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "updateStoragePool")
	if responseError != nil {
		return responseError
	}
	var contentHolder map[string]interface{}
	if err := json.Unmarshal(response, &contentHolder); err != nil {
		log.Printf("Failed to unmarshall response from updateStoragePool: %#v", err)
		return err
	}
	responseHolder := contentHolder["response"].(map[string]interface{})
	anyValueHholder := responseHolder["AnyValue"].(map[string]interface{})
	poolData, err := json.Marshal(anyValueHholder)
	if err != nil {
		panic(err)
	}
	var result storagePool
	if err := json.Unmarshal(poolData, &result); err != nil {
		log.Printf("Failed to unmarshall response from updateStoragePool: %#v", err)
		return err
	}
	err = c.waitForJobCompletion(result.Region, result.Jobs[0].JobID, 1200, 20, false)
	if err != nil {
		return err
	}

	return nil
}
