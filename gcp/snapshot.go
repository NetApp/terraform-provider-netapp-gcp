package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"log"
)

// createSnapshotRequest the users input for creating a Snapshot
type createSnapshotRequest struct {
	Name     string `structs:"name"`
	Region   string `structs:"region"`
	VolumeID string `structs:"volumeId"`
}

// createSnapshotResult the api rsponse for creating a Snapshot
type createSnapshotResult struct {
	Name listSnapshotJobIDResult `json:"response"`
}

// listSnapshotJobIDResult the api rsponse for createSnapshotResult struct creating a Snapshot
type listSnapshotJobIDResult struct {
	JobID listSnapshotIDResult `json:"AnyValue"`
}

// listSnapshotIDResult the api rsponse for listSnapshotJobIDResult struct creating a Snapshot
type listSnapshotIDResult struct {
	SnapshotID string `json:"snapshotId"`
}

// deleteSnapshotRequest the user input for deleteing a Snapshot
type deleteSnapshotRequest struct {
	SnapshotID string `structs:"snapshotId"`
	Region     string `structs:"region"`
	VolumeID   string `structs:"volumeId"`
}

// listSnapshotResult lists the volume for given Snapshot ID
type listSnapshotResult struct {
	SnapshotID     string `json:"snapshotId"`
	LifeCycleState string `json:"lifeCycleState"`
}

// listSnapshotRequest requests the volume for given Snapshot ID and region
type listSnapshotRequest struct {
	SnapshotID string `structs:"snapshotId"`
	Region     string `structs:"region"`
	VolumeID   string `structs:"volumeId"`
}

// updateSnapshotRequest request update name of a snapshot for given Snapshot ID and name
type updateSnapshotRequest struct {
	Name       string `structs:"name"`
	Region     string `structs:"region"`
	VolumeID   string `structs:"volumeId"`
	SnapshotID string `structs:"snapshotId"`
}

func (c *Client) getSnapshotByID(snapshot listSnapshotRequest) (listSnapshotResult, error) {

	baseURL := fmt.Sprintf("%s/Volumes/%s/Snapshots/%s", snapshot.Region, snapshot.VolumeID, snapshot.SnapshotID)

	response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumes request failed")
		return listSnapshotResult{}, err
	}

	var result listSnapshotResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumes")
		return listSnapshotResult{}, err
	}
	if result.LifeCycleState == "deleted" || result.LifeCycleState == "deleting" {
		return listSnapshotResult{}, nil
	}

	return result, nil
}

func (c *Client) createSnapshot(request *createSnapshotRequest) (createSnapshotResult, error) {

	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/Volumes/%s/Snapshots", request.Region, request.VolumeID)
	log.Printf("Parameters: %v", params)

	response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("CreateSnapshot request failed")
		return createSnapshotResult{}, err
	}

	var result createSnapshotResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from CreateSnapshot")
		return createSnapshotResult{}, err
	}

	return result, nil
}

func (c *Client) deleteSnapshot(request deleteSnapshotRequest) error {

	baseURL := fmt.Sprintf("%s/Volumes/%s/Snapshots/%s", request.Region, request.VolumeID, request.SnapshotID)
	_, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("DeleteSnapshot request failed")
		return err
	}
	return nil
}

func (c *Client) updateSnapshot(request updateSnapshotRequest) error {

	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/Volumes/%s/Snapshots/%s", request.Region, request.VolumeID, request.SnapshotID)
	_, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Print("UpdateSnapshot request failed")
		return err
	}

	return nil

}
