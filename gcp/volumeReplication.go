package gcp

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
)

type volumeReplicationRequest struct {
	Bandwidth           string `structs:"bandwidth,omitempty"`
	DestinationVolumeID string `structs:"destinationVolumeUUID,omitempty"`
	EndpointType        string `structs:"endpointType,omitempty"`
	MirrorState         string `structs:"mirrorState,omitempty"`
	Name                string `structs:"name,omitempty"`
	Policy              string `structs:"replicationPolicy,omitempty"`
	Region              string `structs:"region,omitempty"`
	RelationshipStatus  string `structs:"relationshipStatus,omitempty"`
	ReplicationID       string `structs:"volumeReplicationUUID,omitempty"`
	RemoteRegion        string `structs:"remoteRegion,omitempty"`
	Schedule            string `structs:"replicationSchedule,omitempty"`
	SourceVolumeID      string `structs:"sourceVolumeUUID,omitempty"`
}

type volumeReplicationResult struct {
	Bandwidth             string `json:"bandwidth,omitempty"`
	DestinationVolumeID   string `json:"destinationVolumeUUID,omitempty"`
	EndpointType          string `json:"endpointType,omitempty"`
	LifeCycleState        string `json:"lifeCycleState,omitempty"`
	LifeCycleStateDetails string `json:"lifeCycleStateDetails,omitempty"`
	MirrorState           string `json:"mirrorState,omitempty"`
	Name                  string `json:"name,omitempty"`
	Policy                string `json:"replicationPolicy,omitempty"`
	RelationshipStatus    string `json:"relationshipStatus,omitempty"`
	RemoteRegion          string `json:"remoteRegion,omitempty"`
	ReplicationID         string `json:"volumeReplicationUUID,omitempty"`
	Schedule              string `json:"replicationSchedule,omitempty"`
	SourceVolumeID        string `json:"sourceVolumeUUID,omitempty"`
}

func (c *Client) getVolumeReplicationByID(replica volumeReplicationRequest) (volumeReplicationResult, error) {

	baseURL := fmt.Sprintf("%s/VolumeReplications/%s", replica.Region, replica.ReplicationID)

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("getVolumeReplicationByID request failed")
		return volumeReplicationResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "getVolumeReplicationByID")
	if responseError != nil {
		return volumeReplicationResult{}, responseError
	}

	var result volumeReplicationResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumeReplicationByID")
		return volumeReplicationResult{}, err
	}

	if result.LifeCycleState == "deleted" || result.LifeCycleState == "deleting" {
		return volumeReplicationResult{}, nil
	}

	return result, nil
}

func (c *Client) createVolumeReplication(replica *volumeReplicationRequest) (volumeReplicationResult, error) {
	baseURL := fmt.Sprintf("%s/VolumeReplications", replica.Region)

	params := structs.Map(replica)

	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("createVolumeReplication request failed")
		return volumeReplicationResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createVolumeReplication")
	if responseError != nil {
		return volumeReplicationResult{}, responseError
	}

	var f interface{}
	err = json.Unmarshal(response, &f)
	if err != nil {
		return volumeReplicationResult{}, err
	}
	m := f.(map[string]interface{})
	jobs := m["jobs"].([]interface{})
	for _, v := range jobs {
		job := v.(map[string]interface{})
		if job["action"].(string) == "create" {
			err := c.waitForJobCompletion(replica.Region, job["jobId"].(string), 600, 10, false)
			if err != nil {
				return volumeReplicationResult{}, err
			}
		}
	}

	var result volumeReplicationResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createVolumeReplication")
		return volumeReplicationResult{}, err
	}

	return result, nil
}

func (c *Client) breakVolumeReplication(replica *volumeReplicationRequest) error {
	baseURL := fmt.Sprintf("%s/VolumeReplications/%s/Break", replica.Region, replica.ReplicationID)
	statusCode, response, err := c.CallAPIMethod("POST", baseURL, nil)
	if err != nil {
		log.Print("breakVolumeReplication request failed")
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "breakVolumeReplication")
	if responseError != nil {
		return responseError
	}

	var f interface{}
	err = json.Unmarshal(response, &f)
	if err != nil {
		return err
	}
	m := f.(map[string]interface{})
	content := m["response"].(map[string]interface{})
	anyValue := content["AnyValue"].(map[string]interface{})
	jobs := anyValue["jobs"].([]interface{})
	for _, v := range jobs {
		job := v.(map[string]interface{})
		if job["action"].(string) == "break" {
			err := c.waitForJobCompletion(replica.Region, job["jobId"].(string), 600, 10, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Given a jobID and region, wait for the job to finish. All measurments are in seconds.
// if waitUntilCompleted is true, it will not return until the job is done or encounters error.
func (c *Client) waitForJobCompletion(region string, jobID string, timeout int, interval int, waitUntilCompleted bool) error {

	for timeout > 0 || waitUntilCompleted {
		if timeout > 0 {
			timeout -= interval
		}
		time.Sleep(time.Duration(interval) * time.Second)
		jobDetail, err := c.getJobByID(region, jobID)
		if err != nil {
			return err
		}
		if jobDetail.State == "done" {
			return nil
		}
		if jobDetail.State == "error" {
			return fmt.Errorf(jobDetail.StateDetails)
		}
	}
	log.Printf("Job is still onging, return after maximum wait time is reached.")
	return nil
}

func (c *Client) deleteVolumeReplication(replica *volumeReplicationRequest) error {
	baseURL := fmt.Sprintf("%s/VolumeReplications/%s", replica.Region, replica.ReplicationID)

	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("deleteVolumeReplication request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteVolumeReplication")
	if responseError != nil {
		return responseError
	}

	return nil

}

func (c *Client) updateVolumeReplication(replica *volumeReplicationRequest) error {

	baseURL := fmt.Sprintf("%s/VolumeReplications/%s", replica.Region, replica.ReplicationID)

	params := structs.Map(replica)

	statusCode, response, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Print("updateVolumeReplication request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "updateVolumeReplication")
	if responseError != nil {
		return responseError
	}

	return nil
}

func (c *Client) getJobByID(region string, jobID string) (job, error) {

	baseURL := fmt.Sprintf("%s/Jobs/%s", region, jobID)

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("updateVolumeReplication request failed")
		return job{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getJobByID")
	if responseError != nil {
		return job{}, responseError
	}
	var result job
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getJobByID")
		return job{}, err
	}
	log.Printf("get job id: %#v", result)
	return result, nil

}

type job struct {
	State        string `json:"state"`
	StateDetails string `json:"stateDetails"`
}
