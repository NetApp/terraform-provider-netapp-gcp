package gcp

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

// createVolumeBackupRequest the users input for creating a VolumeBackup
type createVolumeBackupRequest struct {
	Name     string `structs:"name"`
	Region   string `structs:"region"`
	VolumeID string `structs:"volumeId"`
}

// createVolumeBackupResult the api response for creating a VolumeBackup
type createVolumeBackupResult struct {
	Name listVolumeBackupJobIDResult `json:"response"`
}

// listVolumeBackupJobIDResult the api response for createVolumeBackupResult struct creating a VolumeBackup
type listVolumeBackupJobIDResult struct {
	JobID listVolumeBackupIDResult `json:"AnyValue"`
}

// listVolumeBackupIDResult the api response for listVolumeBackupJobIDResult struct creating a VolumeBackup
type listVolumeBackupIDResult struct {
	VolumeBackupID string `json:"backupId"`
}

// deleteVolumeBackupRequest the user input for deleteing a VolumeBackup
type deleteVolumeBackupRequest struct {
	VolumeBackupID string `structs:"backupId"`
	Region         string `structs:"region"`
	VolumeID       string `structs:"volumeId"`
}

// listVolumeBackupResult lists the volume for given VolumeBackup ID
type listVolumeBackupResult struct {
	VolumeBackupID string `json:"backupId"`
	LifeCycleState string `json:"lifeCycleState"`
}

// listVolumeBackupRequest requests the volume for given VolumeBackup ID and region
type listVolumeBackupRequest struct {
	VolumeBackupID string `structs:"backupId"`
	Region         string `structs:"region"`
	VolumeID       string `structs:"volumeId"`
}

// updateVolumeBackupRequest request update name of a VolumeBackup for given VolumeBackup ID and name
type updateVolumeBackupRequest struct {
	Name           string `structs:"name"`
	Region         string `structs:"region"`
	VolumeID       string `structs:"volumeId"`
	VolumeBackupID string `structs:"backupId"`
}

func (c *Client) getVolumeBackupByID(VolumeBackup listVolumeBackupRequest) (listVolumeBackupResult, error) {

	baseURL := fmt.Sprintf("%s/Volumes/%s/Backups/%s", VolumeBackup.Region, VolumeBackup.VolumeID, VolumeBackup.VolumeBackupID)

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumeBackup request failed")
		return listVolumeBackupResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "ListVolumeBackup")
	if responseError != nil {
		return listVolumeBackupResult{}, responseError
	}

	var result listVolumeBackupResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumes")
		return listVolumeBackupResult{}, err
	}
	if result.LifeCycleState == "deleted" || result.LifeCycleState == "deleting" {
		return listVolumeBackupResult{}, nil
	}

	return result, nil
}

func (c *Client) createVolumeBackup(request *createVolumeBackupRequest) (createVolumeBackupResult, error) {

	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/Volumes/%s/Backups", request.Region, request.VolumeID)
	log.Printf("Parameters: %v", params)

	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("CreateVolumeBackup request failed")
		return createVolumeBackupResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "CreateVolumeBackup")
	if responseError != nil {
		return createVolumeBackupResult{}, responseError
	}

	var result createVolumeBackupResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from CreateVolumeBackup")
		return createVolumeBackupResult{}, err
	}

	return result, nil
}

func (c *Client) deleteVolumeBackup(request deleteVolumeBackupRequest) error {

	baseURL := fmt.Sprintf("%s/Volumes/%s/Backups/%s", request.Region, request.VolumeID, request.VolumeBackupID)
	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("DeleteVolumeBackup request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "DeleteVolumeBackup")
	if responseError != nil {
		return responseError
	}

	return nil
}
