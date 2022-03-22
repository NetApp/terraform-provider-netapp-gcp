package gcp

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

type kmsConfig struct {
	KeyRing         string `json:"keyRing"`
	KeyName         string `json:"KeyName"`
	KeyRingLocation string `json:"keyRingLocation"`
	ID              string `json:"UUID"`
	KeyProjectID    string `json:"keyProjectID"`
	Network         string `json:"network"`
}

func (c *Client) createKMSConfig(request *kmsConfig) (kmsConfig, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Storage/KmsConfig", request.KeyRingLocation)
	log.Printf("params: %#v", params)
	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("createKMSConfig request failed")
		return kmsConfig{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createKMSConfig")
	if responseError != nil {
		return kmsConfig{}, responseError
	}

	var result kmsConfig
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createKMSConfig")
		return kmsConfig{}, err
	}

	return result, nil
}

func (c *Client) getKMSConfig(request *kmsConfig) (kmsConfig, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Storage/KmsConfig/%s", request.KeyRingLocation, request.ID)
	statusCode, response, err := c.CallAPIMethod("GET", baseURL, params)
	if err != nil {
		log.Print("getKMSConfig request failed")
		return kmsConfig{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getKMSConfig")
	if responseError != nil {
		return kmsConfig{}, responseError
	}
	var result kmsConfig
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getKMSConfig")
		return kmsConfig{}, err
	}

	return result, nil
}

func (c *Client) updateKMSConfig(request *kmsConfig) (kmsConfig, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Storage/KmsConfig/%s", request.KeyRingLocation, request.ID)
	statusCode, response, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Print("updateKMSConfig request failed")
		return kmsConfig{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "updateKMSConfig")
	if responseError != nil {
		return kmsConfig{}, responseError
	}

	var result kmsConfig
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from updateKMSConfig")
		return kmsConfig{}, err
	}

	return result, nil
}

func (c *Client) deleteKMSConfig(request *kmsConfig) (kmsConfig, error) {
	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Storage/KmsConfig/%s", request.KeyRingLocation, request.ID)
	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, params)
	if err != nil {
		log.Print("deleteKMSConfig request failed")
		return kmsConfig{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "deletteKMSConfig")
	if responseError != nil {
		return kmsConfig{}, responseError
	}

	var result kmsConfig
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from deleteKMSConfig")
		return kmsConfig{}, err
	}

	return result, nil
}
