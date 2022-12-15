package gcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

const contextDeadlineExceededErrorMessage = "Post http://cloud-volumes-service.sde.svc.cluster.local/v2/Volumes: context deadline exceeded"

// volumeRequest the users input for creating,requesting,updateing a Volume
// exportPolicy can't set to omitempty because it could be deleted during update.
type volumeRequest struct {
	Name                   string         `structs:"name,omitempty"`
	Region                 string         `structs:"region,omitempty"`
	CreationToken          string         `structs:"creationToken,omitempty"`
	ProtocolTypes          []string       `structs:"protocolTypes,omitempty"`
	Network                string         `structs:"network,omitempty"`
	Size                   int            `structs:"quotaInBytes,omitempty"`
	ServiceLevel           string         `structs:"serviceLevel,omitempty"`
	SnapshotPolicy         snapshotPolicy `structs:"snapshotPolicy,omitempty"`
	ExportPolicy           exportPolicy   `structs:"exportPolicy"`
	VolumeID               string         `structs:"volumeId,omitempty"`
	PoolID                 string         `structs:"poolId,omitempty"`
	Zone                   string         `structs:"zone,omitempty"`
	StorageClass           string         `structs:"storageClass,omitempty"`
	RegionalHA             bool           `structs:"regionalHA,omitempty"`
	SnapshotDirectory      bool           `structs:"snapshotDirectory"`
	UnixPermissions        string         `structs:"unixPermissions,omitempty"`
	SecurityStyle          string         `structs:"securityStyle,omitempty"`
	SharedVpcProjectNumber string
	SmbShareSettings       []string       `structs:"smbShareSettings,omitempty"`
	BillingLabels          []billingLabel `structs:"billingLabels"`
	SnapshotID             string         `structs:"snapshotId"`
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
	PoolID                string         `json:"poolId,omitempty"`
	LifeCycleState        string         `json:"lifeCycleState"`
	LifeCycleStateDetails string         `json:"lifeCycleStateDetails"`
	MountPoints           []mountPoints  `json:"mountPoints,omitempty"`
	Zone                  string         `json:"zone,omitempty"`
	StorageClass          string         `json:"storageClass,omitempty"`
	RegionalHA            bool           `json:"regionalHA,omitempty"`
	TypeDP                bool           `json:"isDataProtection,omitempty"`
	SnapshotDirectory     bool           `json:"snapshotDirectory,omitempty"`
	SmbShareSettings      []string       `json:"smbShareSettings,omitempty"`
	UnixPermissions       string         `json:"unixPermissions,omitempty"`
	SecurityStyle         string         `json:"securityStyle,omitempty"`
	BillingLabels         []billingLabel `json:"billingLabels,omitempty"`
}

type billingLabel struct {
	Key   string `structs:"key,omitempty"`
	Value string `structs:"value,omitempty"`
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

type simpleExportPolicyRule struct {
	Access              string  `structs:"access"`
	AllowedClients      string  `structs:"allowedClients"`
	HasRootAccess       string  `structs:"hasRootAccess"`
	Kerberos5ReadOnly   checked `structs:"kerberos5ReadOnly"`
	Kerberos5ReadWrite  checked `structs:"kerberos5ReadWrite"`
	Kerberos5iReadOnly  checked `structs:"kerberos5iReadOnly"`
	Kerberos5iReadWrite checked `structs:"kerberos5iReadWrite"`
	Kerberos5pReadOnly  checked `structs:"kerberos5pReadOnly"`
	Kerberos5pReadWrite checked `structs:"kerberos5pReadWrite"`
	Nfsv3               checked `structs:"nfsv3"`
	Nfsv4               checked `structs:"nfsv4"`
}

type exportPolicy struct {
	Rules []simpleExportPolicyRule `structs:"rules"`
}

type checked struct {
	Checked bool `structs:"checked"`
}

type mountPoints struct {
	Export       string `structs:"export"`
	Server       string `structs:"server"`
	ProtocolType string `structs:"protocolType"`
}

func (c *Client) getVolumeByID(volume volumeRequest) (volumeResult, error) {
	var baseURL string
	var originalID string = ""

	// terraform import will specify volumeID.
	// Issue is, that volumeID is unqiue per region, but might exist in different regions in same project.
	// For that case, we will support a special ID format for terraform import ADDR ID.
	// ID = <volumeID>:<region>
	s := strings.Split(volume.VolumeID, ":")
	if len(s) == 2 {
		originalID = volume.VolumeID
		volume.VolumeID = s[0]
		volume.Region = s[1]
	}

	if volume.Region == "" {
		// terraform import: ID = <volumeID> and no region specified
		// find all volumes which match VolumeID
		volumes, err := c.filterAllVolumes(func(v volumeResult) bool {
			return v.VolumeID == volume.VolumeID
		})
		if err != nil {
			return volumeResult{}, err
		}

		if len(volumes) == 0 {
			return volumeResult{}, fmt.Errorf("getVolumeByID: No volume found with ID %s", volume.VolumeID)
		}
		if len(volumes) > 1 {
			// return error message which tells user to rerun with ID = <volumeID>:<region> format
			return volumeResult{}, fmt.Errorf(`getVolumeByID: More than one volume found with ID %s. \n
			If this happend while running terraform import, please use \n
			terraform import ADDR ID, with ID using <volumeID>:<region_name> format`, volume.VolumeID)
		}
		if len(volumes) == 1 {
			// we found the right volume to import
			volume.Region = volumes[0].Region
		}
	}

	baseURL = fmt.Sprintf("%s/Volumes/%s", volume.Region, volume.VolumeID)

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		return volumeResult{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getVolumeByID")
	if responseError != nil {
		return volumeResult{}, responseError
	}

	log.Printf("get get get: %#v", bytes.NewBuffer(response).String())
	var result volumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumeByID")
		return volumeResult{}, err
	}

	// volumeID is verified by Terraform. If we use ID = <volumeID>:<region>, we need to revert our ID changes
	if originalID != "" {
		result.VolumeID = originalID
	}
	return result, nil
}

// refactored, but commented, since no code is using it currently
// func (c *Client) getVolumeByRegion(region string) ([]volumeResult, error) {
// 	return c.getVolumes(region)
// }

// Returns volumes of the project. region = "-" for all regions
func (c *Client) getVolumes(region string) ([]volumeResult, error) {

	baseURL := fmt.Sprintf("%s/Volumes", region)
	var result []volumeResult

	statusCode, response, err := c.CallAPIMethod("GET", baseURL, nil)
	if err != nil {
		log.Print("getVolumes request failed")
		return result, err
	}

	responseError := apiResponseChecker(statusCode, response, "getVolumes")
	if responseError != nil {
		return result, responseError
	}

	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumes")
		return result, err
	}
	return result, nil
}

// Filter all volumes of the project by applying a filter function
// Example filter function: func(v volumeResult) bool { return v.VolumeID == "1234-5678-90" }
func (c *Client) filterAllVolumes(f func(volumeResult) bool) ([]volumeResult, error) {
	filteredVolumes := make([]volumeResult, 0)

	vols, err := c.getVolumes("-")
	if err != nil {
		return filteredVolumes, err
	}

	for _, v := range vols {
		if f(v) {
			filteredVolumes = append(filteredVolumes, v)
		}
	}
	return filteredVolumes, nil
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

func (c *Client) createVolume(request *volumeRequest, volType string) (createVolumeResult, error) {

	if request.CreationToken == "" {
		creationToken, err := c.createVolumeCreationToken(*request)
		if err != nil {
			log.Print("CreateVolume request failed")
			return createVolumeResult{}, err
		}
		request.CreationToken = creationToken.CreationToken
	}

	var projectID string
	if request.SharedVpcProjectNumber != "" {
		projectID = request.SharedVpcProjectNumber
	} else {
		projectID = c.GetProjectID()
	}
	request.Network = fmt.Sprintf("projects/%s/global/networks/%s", projectID, request.Network)

	params := structs.Map(request)
	baseURL := fmt.Sprintf("%s/%s", request.Region, volType)
	statusCode, response, err := c.CallAPIMethod("POST", baseURL, params)
	if err != nil {
		return createVolumeResult{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "createVolume")
	if responseError != nil {
		var responseErrorContent apiErrorResponse
		responseContent := bytes.NewBuffer(response).String()
		if err := json.Unmarshal(response, &responseErrorContent); err != nil {
			return createVolumeResult{}, fmt.Errorf(responseContent)
		}
		if responseErrorContent.Code >= 300 || responseErrorContent.Code < 200 {
			spawnJobCreationErrorMessage := fmt.Sprintf("Error creating volume - Cannot spawn additional jobs in %s for this network . Please wait for the ongoing jobs to finish in zone %s and try again", request.Zone, request.Zone)
			log.Printf("* Response error message on createVolume: %v", responseErrorContent.Message)
			if responseErrorContent.Message == spawnJobCreationErrorMessage {
				retries := 10
				for retries > 0 {
					log.Printf("* Retries %v", retries)
					var spawnJobResponseErrorContent apiErrorResponse
					time.Sleep(time.Duration(nextRandomInt(30, 50)) * time.Second)
					statusCode, response, err = c.CallAPIMethod("POST", baseURL, params)
					if err != nil {
						return createVolumeResult{}, err
					}
					responseError = apiResponseChecker(statusCode, response, "createVolume")
					responseContent = bytes.NewBuffer(response).String()
					if err := json.Unmarshal(response, &spawnJobResponseErrorContent); err != nil {
						return createVolumeResult{}, fmt.Errorf(responseContent)
					}
					if spawnJobResponseErrorContent.Code == 0 {
						var result createVolumeResult
						if err := json.Unmarshal(response, &result); err != nil {
							log.Print("Failed to unmarshall response from createVolume")
							return createVolumeResult{}, fmt.Errorf(bytes.NewBuffer(response).String())
						}
						return result, nil
					}
					if spawnJobResponseErrorContent.Message != spawnJobCreationErrorMessage {
						log.Printf("Retry failed spawnJobResponseErrorContent: %v", spawnJobResponseErrorContent.Message)
						return createVolumeResult{}, responseError
					}
					retries--
				}
			} else if responseErrorContent.Message == contextDeadlineExceededErrorMessage {
				retries := 5
				for retries > 0 {
					var contextDeadlineResponseErrorContent apiErrorResponse
					time.Sleep(time.Duration(nextRandomInt(5, 10)) * time.Second)
					statusCode, response, err = c.CallAPIMethod("POST", baseURL, params)
					if err != nil {
						return createVolumeResult{}, err
					}
					responseError = apiResponseChecker(statusCode, response, "createVolume")
					responseContent = bytes.NewBuffer(response).String()
					if err := json.Unmarshal(response, &contextDeadlineResponseErrorContent); err != nil {
						return createVolumeResult{}, fmt.Errorf(responseContent)
					}
					if contextDeadlineResponseErrorContent.Code == 0 {
						var result createVolumeResult
						if err := json.Unmarshal(response, &result); err != nil {
							log.Print("Failed to unmarshall response from createVolume")
							return createVolumeResult{}, fmt.Errorf(bytes.NewBuffer(response).String())
						}
						return result, nil
					}
					if contextDeadlineResponseErrorContent.Message != contextDeadlineExceededErrorMessage {
						return createVolumeResult{}, responseError
					}
					retries--
				}
			} else {
				return createVolumeResult{}, responseError
			}
		}
	}

	var result createVolumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createVolume")
		return createVolumeResult{}, err
	}

	return result, nil
}

func (c *Client) deleteVolume(request volumeRequest) error {
	log.Print("deleteVolume...")
	baseURL := fmt.Sprintf("%s/Volumes/%s", request.Region, request.VolumeID)
	statusCode, response, err := c.CallAPIMethod("DELETE", baseURL, nil)
	if err != nil {
		log.Print("DeleteVolume request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteVolume")
	if responseError != nil {
		var responseErrorContent apiErrorResponse
		responseContent := bytes.NewBuffer(response).String()
		if err := json.Unmarshal(response, &responseErrorContent); err != nil {
			return fmt.Errorf(responseContent)
		}
		if responseErrorContent.Code >= 300 || responseErrorContent.Code < 200 {
			spawnJobDeletionErrorMessage := fmt.Sprintf("Error deleting volume - Cannot spawn additional jobs in %s for this network . Please wait for the ongoing jobs to finish in zone %s and try again", request.Zone, request.Zone)
			log.Printf("* Response error message on deleteVolume: %v", responseErrorContent.Message)
			if responseErrorContent.Message == spawnJobDeletionErrorMessage {
				retries := 10
				for retries > 0 {
					log.Printf("retries %v", retries)
					var deleteJobResponseErrorContent apiErrorResponse
					time.Sleep(time.Duration(nextRandomInt(30, 50)) * time.Second)
					statusCode, response, err = c.CallAPIMethod("DELETE", baseURL, nil)
					if err != nil {
						return err
					}
					responseError = apiResponseChecker(statusCode, response, "deleteVolume")
					responseContent = bytes.NewBuffer(response).String()
					if err := json.Unmarshal(response, &deleteJobResponseErrorContent); err != nil {
						return fmt.Errorf(responseContent)
					}
					if deleteJobResponseErrorContent.Code == 0 {
						var result createVolumeResult
						if err := json.Unmarshal(response, &result); err != nil {
							log.Print("Failed to unmarshall response from createVolume")
							return fmt.Errorf(bytes.NewBuffer(response).String())
						}
						return nil
					}
					if deleteJobResponseErrorContent.Message != spawnJobDeletionErrorMessage {
						log.Printf("Retry failed deleteJobResponseErrorContent: %v", deleteJobResponseErrorContent.Message)
						return responseError
					}
					retries--
				}

			} else {
				return responseError
			}
		}
	}

	var result apiErrorResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from deleteVolume")
		return fmt.Errorf(bytes.NewBuffer(response).String())
	}

	return nil
}

func (c *Client) createVolumeCreationToken(request volumeRequest) (volumeResult, error) {
	params := structs.Map(request)

	baseURL := fmt.Sprintf("%s/VolumeCreationToken", request.Region)
	log.Printf("Parameters: %v", params)
	statusCode, response, err := c.CallAPIMethod("GET", baseURL, params)
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
	snapshotPolicy := snapshotPolicy{}

	if v, ok := data["enabled"]; ok {
		snapshotPolicy.Enabled = v.(bool)
	}

	if v, ok := data["daily_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			dailySchedule := v.([]interface{})[0].(map[string]interface{})
			if hour, ok := dailySchedule["hour"]; ok {
				snapshotPolicy.DailySchedule.Hour = hour.(int)
			}
			if minute, ok := dailySchedule["minute"]; ok {
				snapshotPolicy.DailySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := dailySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.DailySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["hourly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			hourlySchedule := v.([]interface{})[0].(map[string]interface{})
			if minute, ok := hourlySchedule["minute"]; ok {
				snapshotPolicy.HourlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := hourlySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.HourlySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["monthly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			monthlySchedule := v.([]interface{})[0].(map[string]interface{})
			if daysOfMonth, ok := monthlySchedule["days_of_month"]; ok {
				snapshotPolicy.MonthlySchedule.DaysOfMonth = daysOfMonth.(string)
			}
			if hour, ok := monthlySchedule["hour"]; ok {
				snapshotPolicy.MonthlySchedule.Hour = hour.(int)
			}
			if minute, ok := monthlySchedule["minute"]; ok {
				snapshotPolicy.MonthlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := monthlySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.MonthlySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["weekly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			weeklySchedule := v.([]interface{})[0].(map[string]interface{})
			if day, ok := weeklySchedule["day"]; ok {
				snapshotPolicy.WeeklySchedule.Day = day.(string)
			}
			if hour, ok := weeklySchedule["hour"]; ok {
				snapshotPolicy.WeeklySchedule.Hour = hour.(int)
			}
			if minute, ok := weeklySchedule["minute"]; ok {
				snapshotPolicy.WeeklySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := weeklySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.WeeklySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	return snapshotPolicy
}

// flattenExportPolicy converts exportPolicy struct to []map[string]interface{}
func flattenExportPolicy(v exportPolicy) interface{} {
	exportPolicyRules := v.Rules
	rules := make([]map[string]interface{}, 0, len(exportPolicyRules))
	for _, exportPolicyRule := range exportPolicyRules {
		ruleMap := make(map[string]interface{})
		ruleMap["access"] = exportPolicyRule.Access
		ruleMap["allowed_clients"] = exportPolicyRule.AllowedClients
		ruleMap["has_root_access"] = exportPolicyRule.HasRootAccess
		ruleMap["kerberos5_readonly"] = exportPolicyRule.Kerberos5ReadOnly.Checked
		ruleMap["kerberos5_readwrite"] = exportPolicyRule.Kerberos5ReadWrite.Checked
		ruleMap["kerberos5i_readonly"] = exportPolicyRule.Kerberos5iReadOnly.Checked
		ruleMap["kerberos5i_readwrite"] = exportPolicyRule.Kerberos5iReadWrite.Checked
		ruleMap["kerberos5p_readonly"] = exportPolicyRule.Kerberos5pReadOnly.Checked
		ruleMap["kerberos5p_readwrite"] = exportPolicyRule.Kerberos5pReadWrite.Checked
		nfsv3Config := make(map[string]interface{})
		nfsv4Config := make(map[string]interface{})
		nfsv3Config["checked"] = exportPolicyRule.Nfsv3.Checked
		nfsv4Config["checked"] = exportPolicyRule.Nfsv4.Checked
		nfsv3 := make([]map[string]interface{}, 1)
		nfsv4 := make([]map[string]interface{}, 1)
		nfsv3[0] = make(map[string]interface{})
		nfsv4[0] = make(map[string]interface{})
		nfsv3[0] = nfsv3Config
		nfsv4[0] = nfsv4Config
		ruleMap["nfsv3"] = nfsv3
		ruleMap["nfsv4"] = nfsv4
		rules = append(rules, ruleMap)
	}
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["rule"] = rules
	return result
}

// expandExportPolicy converts set to exportPolicy struct
func expandExportPolicy(set *schema.Set) (exportPolicy, error) {
	exportPolicyObj := exportPolicy{}

	for _, v := range set.List() {
		rules := v.(map[string]interface{})
		ruleSet := rules["rule"].([]interface{})
		ruleConfigs := make([]simpleExportPolicyRule, 0, len(ruleSet))
		for _, x := range ruleSet {
			exportPolicyRule := simpleExportPolicyRule{}
			ruleConfig := x.(map[string]interface{})
			exportPolicyRule.Access = ruleConfig["access"].(string)
			exportPolicyRule.AllowedClients = ruleConfig["allowed_clients"].(string)
			exportPolicyRule.HasRootAccess = ruleConfig["has_root_access"].(string)
			exportPolicyRule.Kerberos5ReadOnly.Checked = ruleConfig["kerberos5_readonly"].(bool)
			exportPolicyRule.Kerberos5ReadWrite.Checked = ruleConfig["kerberos5_readwrite"].(bool)
			exportPolicyRule.Kerberos5iReadOnly.Checked = ruleConfig["kerberos5i_readonly"].(bool)
			exportPolicyRule.Kerberos5iReadWrite.Checked = ruleConfig["kerberos5i_readwrite"].(bool)
			exportPolicyRule.Kerberos5pReadOnly.Checked = ruleConfig["kerberos5p_readonly"].(bool)
			exportPolicyRule.Kerberos5pReadWrite.Checked = ruleConfig["kerberos5p_readwrite"].(bool)
			nfsv3Set := ruleConfig["nfsv3"].(*schema.Set)
			nfsv4Set := ruleConfig["nfsv4"].(*schema.Set)
			for _, y := range nfsv3Set.List() {
				nfsv3Config := y.(map[string]interface{})
				exportPolicyRule.Nfsv3.Checked = nfsv3Config["checked"].(bool)
			}
			for _, z := range nfsv4Set.List() {
				nfsv4Config := z.(map[string]interface{})
				exportPolicyRule.Nfsv4.Checked = nfsv4Config["checked"].(bool)
			}
			if !exportPolicyRule.Nfsv3.Checked && !exportPolicyRule.Nfsv4.Checked {
				return exportPolicy{}, fmt.Errorf("At least one of nfsv3 or nfsv4 needs to be true in protocol type of the export policy rule")
			}
			ruleConfigs = append(ruleConfigs, exportPolicyRule)
		}
		exportPolicyObj.Rules = ruleConfigs
	}
	return exportPolicyObj, nil
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

func flattenBillingLabel(v []billingLabel) interface{} {
	labels := make([]map[string]interface{}, 0, len(v))
	for _, l := range v {
		labelMap := make(map[string]interface{})
		labelMap["key"] = l.Key
		labelMap["value"] = l.Value
		labels = append(labels, labelMap)
	}
	return labels
}

func expandBillingLabel(set *schema.Set) []billingLabel {
	var billingLabels []billingLabel
	for _, v := range set.List() {
		billingLabel := billingLabel{}
		blabel := v.(map[string]interface{})
		billingLabel.Key = blabel["key"].(string)
		billingLabel.Value = blabel["value"].(string)
		billingLabels = append(billingLabels, billingLabel)
	}
	return billingLabels
}
