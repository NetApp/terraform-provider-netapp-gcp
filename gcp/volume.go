package gcp

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

// createVolumeRequest the users input for creating a Volume
type createVolumeRequest struct {
	Name           string         `structs:"name"`
	Region         string         `structs:"region"`
	CreationToken  string         `structs:"creationToken"`
	ProtocolTypes  []string       `structs:"protocolTypes"`
	Network        string         `structs:"network"`
	Size           int            `structs:"quotaInBytes"`
	ServiceLevel   string         `structs:"serviceLevel"`
	SnapshotPolicy snapshotPolicy `structs:"snapshotPolicy"`
	ExportPolicy   exportPolicy   `structs:"exportPolicy"`
}

// createVolumeResult the api response for creating a volume
type createVolumeResult struct {
	Name    listVolumeJobIDResult `json:"response"`
	Code    int                   `json:"code"`
	Message string                `json:"message"`
}

// listVolumeResult requests the volume for given volume ID and region
type listVolumesRequest struct {
	VolumeID string `structs:"volumeId"`
	Region   string `structs:"region"`
}

// listVolumeResult lists the volume for given volume ID
type listVolumeResult struct {
	VolumeID              string       `json:"volumeId"`
	VolumeName            string       `json:"name"`
	CreationToken         string       `json:"creationToken"`
	LifeCycleState        string       `json:"lifeCycleState"`
	LifeCycleStateDetails string       `json:"lifeCycleStateDetails"`
	MountPoints           []mountPoint `json:"mountPoints"`
}

// listVolumesByNameRequest requests the volume for given volume ID and region
type listVolumesByNameRequest struct {
	Region        string `structs:"region"`
	VolumeID      string `json:"volumeId"`
	VolumeName    string `json:"name"`
	CreationToken string `json:"creationToken"`
}

type updateVolumeRequest struct {
	Name           string         `structs:"name"`
	Region         string         `structs:"region"`
	ProtocolTypes  []string       `structs:"protocolTypes"`
	Network        string         `structs:"network"`
	Size           int            `structs:"quotaInBytes"`
	ServiceLevel   string         `structs:"serviceLevel"`
	SnapshotPolicy snapshotPolicy `structs:"snapshotPolicy"`
	ExportPolicy   exportPolicy   `structs:"exportPolicy"`
	VolumeID       string         `structs:"volumeId"`
}

// listVolumeJobIDResult the api response for createVolumeResult struct creating a volume
type listVolumeJobIDResult struct {
	JobID listVolumeIDResult `json:"AnyValue"`
}

// listVolumeIDResult the api response for listVolumeJobIDResult struct creating a volume
type listVolumeIDResult struct {
	VolID string `json:"volumeId"`
}

// createVolumeCreationTokenResult the api results for creating a volume
type createVolumeCreationTokenResult struct {
	CreationTokenName string `json:"creationToken"`
}

// deleteVolumeRequest the user input for deleteing a volume
type deleteVolumeRequest struct {
	VolumeID string `structs:"volumeId"`
	Region   string `structs:"region"`
}

type mountPoint struct {
	Export       string `structs:"export"`
	ExportFull   string `structs:"exportFull"`
	ProtocolType string `structs:"protocolType"`
	Server       string `structs:"server"`
}

type snapshotPolicy struct {
	Enabled         bool            `structs:"enabled"`
	DailySchedule   dailySchedule   `structs:"dailySchedule"`
	HourlySchedule  hourlySchedule  `structs:"hourlySchedule"`
	MonthlySchedule monthlySchedule `structs:"monthlySchedule"`
	WeeklySchedule  weeklySchedule  `structs:"weeklySchedule"`
}

type dailySchedule struct {
	Hour            int `structs:"hour"`
	Minute          int `structs:"minute"`
	SnapshotsToKeep int `structs:"snapshotsToKeep"`
}

type hourlySchedule struct {
	Minute          int `structs:"minute"`
	SnapshotsToKeep int `structs:"snapshotsToKeep"`
}

type monthlySchedule struct {
	DaysOfMonth     string `structs:"daysOfMonth"`
	Hour            int    `structs:"hour"`
	Minute          int    `structs:"minute"`
	SnapshotsToKeep int    `structs:"snapshotsToKeep"`
}

type weeklySchedule struct {
	Day             string `structs:"day"`
	Hour            int    `structs:"hour"`
	Minute          int    `structs:"minute"`
	SnapshotsToKeep int    `structs:"snapshotsToKeep"`
}

type apiResponseCodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type exportPolicyRule struct {
	Access         string `structs:"access"`
	AllowedClients string `structs:"allowedClients"`
	Nfsv3          nfs    `structs:"nfsv3"`
	Nfsv4          nfs    `structs:"nfsv4"`
}

type exportPolicy struct {
	Rules []exportPolicyRule `structs:"rules"`
}

type nfs struct {
	Checked bool `structs:"checked"`
}

type simpleExportPolicyRule struct {
	SimpleExportPolicyRule exportPolicyRule `structs:"SimpleExportPolicyRule"`
}

func (c *Client) getVolumeByID(volume listVolumesRequest) (listVolumeResult, error) {

	baseURL := fmt.Sprintf("%s/Volumes/%s", volume.Region, volume.VolumeID)

	response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumes request failed")
		return listVolumeResult{}, err
	}

	var result listVolumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumes")
		return listVolumeResult{}, err
	}

	if result.LifeCycleState == "deleted" || result.LifeCycleState == "deleting" {
		return listVolumeResult{}, nil
	}

	return result, nil
}

func (c *Client) getVolumeByNameOrCreationToken(volume listVolumesByNameRequest) (listVolumeResult, error) {

	if volume.VolumeName == "" && volume.CreationToken == "" {
		return listVolumeResult{}, fmt.Errorf("Either CreationToken or volume name or both are required")
	}

	baseURL := fmt.Sprintf("%s/Volumes", volume.Region)

	response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumesByName request failed")
		return listVolumeResult{}, err
	}

	var result []listVolumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumes")
		return listVolumeResult{}, err
	}

	var count = 0
	var resultVolume listVolumeResult
	for _, eachVolume := range result {
		if volume.CreationToken != "" && eachVolume.CreationToken == volume.CreationToken {
			if volume.VolumeName != "" && eachVolume.VolumeName == volume.VolumeName {
				return eachVolume, nil
			} else if volume.VolumeName != "" && eachVolume.VolumeName != volume.VolumeName {
				return listVolumeResult{}, fmt.Errorf("Given CreationToken does not match with given volume name : %v", volume.VolumeName)
			}
			return eachVolume, nil
		} else if volume.CreationToken == "" && volume.VolumeName != "" && eachVolume.VolumeName == volume.VolumeName {
			count = count + 1
			resultVolume = eachVolume
		}
	}

	if volume.CreationToken != "" {
		return listVolumeResult{}, fmt.Errorf("Given CreationToken does not exist : %v", volume.CreationToken)
	}
	if count > 1 {
		return listVolumeResult{}, fmt.Errorf("Found more than one volume : %v", volume.VolumeName)
	} else if count == 0 {
		return listVolumeResult{}, fmt.Errorf("No volume found for : %v", volume.VolumeName)
	}

	return resultVolume, nil
}

func (c *Client) createVolume(request *createVolumeRequest) (createVolumeResult, error) {
	creationToken, err := c.createVolumeCreationToken(*request)
	if err != nil {
		log.Print("CreateVolume request failed")
		return createVolumeResult{}, err
	}

	request.CreationToken = creationToken.CreationTokenName
	projectID := c.GetProjectID()
	request.Network = fmt.Sprintf("projects/%s/global/networks/%s", projectID, request.Network)

	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/Volumes", request.Region)
	log.Printf("Parameters: %v", params)

	response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("CreateVolume request failed")
		return createVolumeResult{}, err
	}

	var result createVolumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from CreateVolume")
		return createVolumeResult{}, err
	}

	//size too small or too large return 500 status code and a error message.
	if (result.Code != 0 && result.Code != 200) || (result.Message != "") {
		return createVolumeResult{}, fmt.Errorf("code: %d, message: %s", result.Code, result.Message)
	}

	return result, nil
}

func (c *Client) deleteVolume(request deleteVolumeRequest) error {

	baseURL := fmt.Sprintf("%s/Volumes/%s", request.Region, request.VolumeID)
	response, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("DeleteVolume request failed")
		return err
	}

	var result apiResponseCodeMessage
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from CreationToken")
		return err
	}
	if (result.Code != 0 && result.Code != 200) || (result.Message != "") {
		return fmt.Errorf("code: %d, message: %s", result.Code, result.Message)
	}

	return nil
}

func (c *Client) createVolumeCreationToken(request createVolumeRequest) (createVolumeCreationTokenResult, error) {
	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/VolumeCreationToken", request.Region)
	log.Printf("Parameters: %v", params)

	response, err := c.CallAPIMethod("", baseURL, params)
	if err != nil {
		log.Print("CreationToken request failed")
		return createVolumeCreationTokenResult{}, err
	}

	var result createVolumeCreationTokenResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from CreationToken")
		return createVolumeCreationTokenResult{}, err
	}
	return result, nil
}

func (c *Client) updateVolume(request updateVolumeRequest) error {
	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/Volumes/%s", request.Region, request.VolumeID)

	response, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Print("updateVolume request failed")
		return err
	}
	var result apiResponseCodeMessage
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from updateVolume")
		return err
	}
	if (result.Code != 0 && result.Code != 200) || (result.Message != "") {
		return fmt.Errorf("code: %d, message: %s", result.Code, result.Message)
	}

	return nil
}

// SetProjectID for the client to use for requests to the GCP API
func (c *Client) SetProjectID(project string) {
	c.Project = project
}

// GetProjectID returns the API version that will be used for GCP API requests
func (c *Client) GetProjectID() string {
	return c.Project
}

func expandSnapshotPolicy(data map[string]interface{}) snapshotPolicy {
	snapshot_policy := snapshotPolicy{}

	if v, ok := data["enabled"]; ok {
		snapshot_policy.Enabled = v.(bool)
	}

	if v, ok := data["daily_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			daily_schedule := v.([]interface{})[0].(map[string]interface{})
			if hour, ok := daily_schedule["hour"]; ok {
				snapshot_policy.DailySchedule.Hour = hour.(int)
			}
			if minute, ok := daily_schedule["minute"]; ok {
				snapshot_policy.DailySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := daily_schedule["snapshots_to_keep"]; ok {
				snapshot_policy.DailySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["hourly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			hourly_schedule := v.([]interface{})[0].(map[string]interface{})
			if minute, ok := hourly_schedule["minute"]; ok {
				snapshot_policy.HourlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := hourly_schedule["snapshots_to_keep"]; ok {
				snapshot_policy.HourlySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["montly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			montly_schedule := v.([]interface{})[0].(map[string]interface{})
			if days_of_month, ok := montly_schedule["days_of_month"]; ok {
				snapshot_policy.MonthlySchedule.DaysOfMonth = days_of_month.(string)
			}
			if hour, ok := montly_schedule["hour"]; ok {
				snapshot_policy.MonthlySchedule.Hour = hour.(int)
			}
			if minute, ok := montly_schedule["minute"]; ok {
				snapshot_policy.MonthlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := montly_schedule["snapshots_to_keep"]; ok {
				snapshot_policy.MonthlySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["weekly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			weekly_schedule := v.([]interface{})[0].(map[string]interface{})
			if day, ok := weekly_schedule["day"]; ok {
				snapshot_policy.WeeklySchedule.Day = day.(string)
			}
			if hour, ok := weekly_schedule["hour"]; ok {
				snapshot_policy.WeeklySchedule.Hour = hour.(int)
			}
			if minute, ok := weekly_schedule["minute"]; ok {
				snapshot_policy.WeeklySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := weekly_schedule["snapshots_to_keep"]; ok {
				snapshot_policy.WeeklySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	return snapshot_policy
}

func expandExportPolicy(data map[string]interface{}) exportPolicy {
	export_policy := exportPolicy{}

	if v, ok := data["rule"]; ok {

		for _, value := range v.([]interface{}) {
			rule := exportPolicyRule{}
			rule_map := value.(map[string]interface{})
			if access := rule_map["access"]; access != "" {
				rule.Access = access.(string)
			}
			if allowedClients := rule_map["allowed_clients"]; allowedClients != "" {
				rule.AllowedClients = allowedClients.(string)
			}
			nfs := nfs{}
			nfsv3 := rule_map["nfsv3"]
			if len(nfsv3.([]interface{})) > 0 {
				nfsv3 := nfsv3.([]interface{})[0].(map[string]interface{})
				if checked, ok := nfsv3["checked"]; ok {
					nfs.Checked = checked.(bool)
				}
				rule.Nfsv3 = nfs
			}
			nfsv4 := rule_map["nfsv4"]
			if len(nfsv4.([]interface{})) > 0 {
				nfsv4 := nfsv4.([]interface{})[0].(map[string]interface{})
				if checked, ok := nfsv4["checked"]; ok {
					nfs.Checked = checked.(bool)
				}
				rule.Nfsv4 = nfs
			}
			export_policy.Rules = append(export_policy.Rules, rule)
		}

		return export_policy
	}

	return exportPolicy{}
}
