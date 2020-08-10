package gcp

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

// volumeRequest the users input for creating,requesting,updateing a Volume
// exportPolicy can't set to omitempty because it could be deleted during update.
type volumeRequest struct {
	Name                      string         `structs:"name,omitempty"`
	Region                    string         `structs:"region,omitempty"`
	CreationToken             string         `structs:"creationToken,omitempty"`
	ProtocolTypes             []string       `structs:"protocolTypes,omitempty"`
	Network                   string         `structs:"network,omitempty"`
	Size                      int            `structs:"quotaInBytes,omitempty"`
	ServiceLevel              string         `structs:"serviceLevel,omitempty"`
	SnapshotPolicy            snapshotPolicy `structs:"snapshotPolicy,omitempty"`
	ExportPolicy              exportPolicy   `structs:"exportPolicy"`
	VolumeID                  string         `structs:"volumeId,omitempty"`
	Shared_vpc_project_number string
}

// volumeRequest retrieves the volume attributes from API and convert to struct
type volumeResult struct {
	Name                  string         `json:"name,omitempty"`
	Region                string         `json:"region,omitempty"`
	CreationToken         string         `json:"creationToken,omitempty"`
	ProtocolTypes         []string       `json:"protocolTypes,omitempty"`
	Network               string         `json:"network,omitempty"`
	Size                  int            `json:"quotaInBytes,omitempty"`
	ServiceLevel          string         `json:"serviceLevel,omitempty"`
	SnapshotPolicy        snapshotPolicy `json:"snapshotPolicy,omitempty"`
	ExportPolicy          exportPolicy   `json:"exportPolicy,omitempty"`
	VolumeID              string         `json:"volumeId,omitempty"`
	LifeCycleState        string         `json:"lifeCycleState"`
	LifeCycleStateDetails string         `json:"lifeCycleStateDetails"`
	MountPoints           []mountPoints  `json:"mountPoints,omitempty"`
}

// createVolumeResult the api response for creating a volume
type createVolumeResult struct {
	Name    listVolumeJobIDResult `json:"response"`
	Code    int                   `json:"code"`
	Message string                `json:"message"`
}

// listVolumeJobIDResult the api response for createVolumeResult struct creating a volume
type listVolumeJobIDResult struct {
	JobID listVolumeIDResult `json:"AnyValue"`
}

// listVolumeIDResult the api response for listVolumeJobIDResult struct creating a volume
type listVolumeIDResult struct {
	VolID string `json:"volumeId"`
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

type mountPoints struct {
	Export       string `structs:"export"`
	Server       string `structs:"server"`
	ProtocolType string `structs:"protocolType"`
}

func (c *Client) getVolumeByID(volume volumeRequest) (volumeResult, error) {

	baseURL := fmt.Sprintf("%s/Volumes/%s", volume.Region, volume.VolumeID)

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumes request failed")
		return volumeResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "getVolumeByID")
	if responseError != nil {
		return volumeResult{}, responseError
	}

	var result volumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumeByID")
		return volumeResult{}, err
	}

	if result.LifeCycleState == "deleted" || result.LifeCycleState == "deleting" {
		return volumeResult{}, nil
	}

	return result, nil
}

func (c *Client) getVolumeByRegion(region string) ([]volumeResult, error) {

	baseURL := fmt.Sprintf("%s/Volumes", region)
	var volumes []volumeResult

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumes request failed")
		return volumes, err
	}

	responseError := apiResponseChecker(statusCode, response, "getVolumeByRegion")
	if responseError != nil {
		return volumes, responseError
	}

	if err := json.Unmarshal(response, &volumes); err != nil {
		log.Print("Failed to unmarshall response from getVolumeByRegion")
		return volumes, err
	}
	return volumes, nil
}

func (c *Client) getVolumeByNameOrCreationToken(volume volumeRequest) (volumeResult, error) {

	if volume.Name == "" && volume.CreationToken == "" {
		return volumeResult{}, fmt.Errorf("Either CreationToken or volume name or both are required")
	}

	baseURL := fmt.Sprintf("%s/Volumes", volume.Region)

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("ListVolumesByName request failed")
		return volumeResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "getVolumeByNameOrCreationToken")
	if responseError != nil {
		return volumeResult{}, responseError
	}

	var result []volumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumeByNameOrCreationToken")
		return volumeResult{}, err
	}

	var count = 0
	var resultVolume volumeResult
	for _, eachVolume := range result {
		if volume.CreationToken != "" && eachVolume.CreationToken == volume.CreationToken {
			if volume.Name != "" && eachVolume.Name == volume.Name {
				return eachVolume, nil
			} else if volume.Name != "" && eachVolume.Name != volume.Name {
				return volumeResult{}, fmt.Errorf("Given CreationToken does not match with given volume name : %v", volume.Name)
			}
			return eachVolume, nil
		} else if volume.CreationToken == "" && volume.Name != "" && eachVolume.Name == volume.Name {
			count = count + 1
			resultVolume = eachVolume
		}
	}

	if volume.CreationToken != "" {
		return volumeResult{}, fmt.Errorf("Given CreationToken does not exist : %v", volume.CreationToken)
	}
	if count > 1 {
		return volumeResult{}, fmt.Errorf("Found more than one volume : %v", volume.Name)
	} else if count == 0 {
		return volumeResult{}, fmt.Errorf("No volume found for : %v", volume.Name)
	}

	return resultVolume, nil
}

func (c *Client) createVolume(request *volumeRequest) (createVolumeResult, error) {

	volumes, err := c.getVolumeByRegion(request.Region)
	if err != nil {
		return createVolumeResult{}, err
	}

	if request.CreationToken == "" {
		creationToken, err := c.createVolumeCreationToken(*request)
		if err != nil {
			log.Print("CreateVolume request failed")
			return createVolumeResult{}, err
		}
		request.CreationToken = creationToken.CreationToken
	}

	var projectID string
	if request.Shared_vpc_project_number != "" {
		projectID = request.Shared_vpc_project_number
	} else {
		projectID = c.GetProjectID()
	}
	request.Network = fmt.Sprintf("projects/%s/global/networks/%s", projectID, request.Network)

	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/Volumes", request.Region)
	log.Printf("Parameters: %v", params)

	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		log.Print("CreateVolume request failed")
		return createVolumeResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createVolume")
	if responseError != nil {
		return createVolumeResult{}, responseError
	}

	var result createVolumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createVolume")
		return createVolumeResult{}, err
	}

	var volume volumeResult
	volume, err = c.getVolumeByID(volumeRequest{Region: request.Region, VolumeID: result.Name.JobID.VolID})
	if err != nil {
		return createVolumeResult{}, err
	}

	// For each region, the first volume creation will take a long time, usually around 3-4 minutes. Because it’s spinning up an svm.
	// Wait until the volume created if it's the first volume of the region.
	if len(volumes) == 0 {
		wait_seconds := 300
		for wait_seconds > 0 && volume.LifeCycleState == "creating" {
			time.Sleep(10)
			volume, err = c.getVolumeByID(volumeRequest{VolumeID: result.Name.JobID.VolID})
			wait_seconds = wait_seconds - 10
		}
	}

	return result, nil
}

func (c *Client) deleteVolume(request volumeRequest) error {

	baseURL := fmt.Sprintf("%s/Volumes/%s", request.Region, request.VolumeID)
	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("DeleteVolume request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteVolume")
	if responseError != nil {
		return responseError
	}

	var result apiResponseCodeMessage
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from deleteVolume")
		return err
	}

	return nil
}

func (c *Client) createVolumeCreationToken(request volumeRequest) (volumeResult, error) {
	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/VolumeCreationToken", request.Region)
	log.Printf("Parameters: %v", params)

	statusCode, response, err := c.CallAPIMethod("", baseURL, params)
	if err != nil {
		log.Print("CreationToken request failed")
		return volumeResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createVolumeCreationToken")
	if responseError != nil {
		return volumeResult{}, responseError
	}

	var result volumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createVolumeCreationToken")
		return volumeResult{}, err
	}
	return result, nil
}

func (c *Client) updateVolume(request volumeRequest) error {
	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/Volumes/%s", request.Region, request.VolumeID)

	statusCode, response, err := c.CallAPIMethod("PUT", baseURL, params)
	if err != nil {
		log.Print("updateVolume request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "updateVolume")
	if responseError != nil {
		return responseError
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

// expandSnapshotPolicy converts map to snapshotPolicy struct
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
	if v, ok := data["monthly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			monthly_schedule := v.([]interface{})[0].(map[string]interface{})
			if days_of_month, ok := monthly_schedule["days_of_month"]; ok {
				snapshot_policy.MonthlySchedule.DaysOfMonth = days_of_month.(string)
			}
			if hour, ok := monthly_schedule["hour"]; ok {
				snapshot_policy.MonthlySchedule.Hour = hour.(int)
			}
			if minute, ok := monthly_schedule["minute"]; ok {
				snapshot_policy.MonthlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := monthly_schedule["snapshots_to_keep"]; ok {
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

// flattenExportPolicy converts exportPolicy struct to []map[string]interface{}
func flattenExportPolicy(v exportPolicy) interface{} {
	export_policy_rules := v.Rules
	rules := make([]map[string]interface{}, 0, len(export_policy_rules))
	for _, export_policy_rule := range export_policy_rules {
		rule_map := make(map[string]interface{})
		rule_map["access"] = export_policy_rule.Access
		rule_map["allowed_clients"] = export_policy_rule.AllowedClients
		nfsv3Config := make(map[string]interface{})
		nfsv4Config := make(map[string]interface{})
		nfsv3Config["checked"] = export_policy_rule.Nfsv3.Checked
		nfsv4Config["checked"] = export_policy_rule.Nfsv4.Checked
		nfsv3 := make([]map[string]interface{}, 1)
		nfsv4 := make([]map[string]interface{}, 1)
		nfsv3[0] = make(map[string]interface{})
		nfsv4[0] = make(map[string]interface{})
		nfsv3[0] = nfsv3Config
		nfsv4[0] = nfsv4Config
		rule_map["nfsv3"] = nfsv3
		rule_map["nfsv4"] = nfsv4
		rules = append(rules, rule_map)
	}
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["rule"] = rules
	return result
}

// expandExportPolicy converts set to exportPolicy struct
func expandExportPolicy(set *schema.Set) exportPolicy {
	export_policy := exportPolicy{}

	for _, v := range set.List() {
		rules := v.(map[string]interface{})
		ruleSet := rules["rule"].(*schema.Set)
		ruleConfigs := make([]exportPolicyRule, 0, ruleSet.Len())
		for _, x := range ruleSet.List() {
			export_policy_rule := exportPolicyRule{}
			ruleConfig := x.(map[string]interface{})
			export_policy_rule.Access = ruleConfig["access"].(string)
			export_policy_rule.AllowedClients = ruleConfig["allowed_clients"].(string)
			nfsv3Set := ruleConfig["nfsv3"].(*schema.Set)
			nfsv4Set := ruleConfig["nfsv4"].(*schema.Set)
			for _, y := range nfsv3Set.List() {
				nfsv3Config := y.(map[string]interface{})
				export_policy_rule.Nfsv3.Checked = nfsv3Config["checked"].(bool)
			}
			for _, z := range nfsv4Set.List() {
				nfsv4Config := z.(map[string]interface{})
				export_policy_rule.Nfsv4.Checked = nfsv4Config["checked"].(bool)
			}
			ruleConfigs = append(ruleConfigs, export_policy_rule)
		}
		export_policy.Rules = ruleConfigs
	}
	return export_policy
}

// flattenSnapshotPolicy converts snapshotPolicy struct to []map[string]interface{}
func flattenSnapshotPolicy(v snapshotPolicy) interface{} {
	flattened := make([]map[string]interface{}, 1)
	sp := make(map[string]interface{})
	sp["enabled"] = v.Enabled
	hourly := make([]map[string]interface{}, 1)
	hourly[0] = make(map[string]interface{})
	hourly[0]["minute"] = v.HourlySchedule.Minute
	hourly[0]["snapshots_to_keep"] = v.HourlySchedule.SnapshotsToKeep
	daily := make([]map[string]interface{}, 1)
	daily[0] = make(map[string]interface{})
	daily[0]["hour"] = v.DailySchedule.Hour
	daily[0]["minute"] = v.DailySchedule.Minute
	daily[0]["snapshots_to_keep"] = v.DailySchedule.SnapshotsToKeep
	monthly := make([]map[string]interface{}, 1)
	monthly[0] = make(map[string]interface{})
	monthly[0]["days_of_month"] = v.MonthlySchedule.DaysOfMonth
	monthly[0]["hour"] = v.MonthlySchedule.Hour
	monthly[0]["minute"] = v.MonthlySchedule.Minute
	monthly[0]["snapshots_to_keep"] = v.MonthlySchedule.SnapshotsToKeep
	weekly := make([]map[string]interface{}, 1)
	weekly[0] = make(map[string]interface{})
	weekly[0]["day"] = v.WeeklySchedule.Day
	weekly[0]["hour"] = v.WeeklySchedule.Hour
	weekly[0]["minute"] = v.WeeklySchedule.Minute
	weekly[0]["snapshots_to_keep"] = v.WeeklySchedule.SnapshotsToKeep
	sp["daily_schedule"] = daily
	sp["hourly_schedule"] = hourly
	sp["weekly_schedule"] = weekly
	sp["monthly_schedule"] = monthly
	flattened[0] = sp
	return flattened
}

func flattenMountPoints(v []mountPoints) interface{} {
	mps := make([]map[string]interface{}, 0, len(v))
	for _, mountpoint := range v {
		mpmap := make(map[string]interface{})
		mpmap["export"] = mountpoint.Export
		mpmap["server"] = mountpoint.Server
		mpmap["protocol_type"] = mountpoint.ProtocolType
		mps = append(mps, mpmap)
	}
	return mps
}
