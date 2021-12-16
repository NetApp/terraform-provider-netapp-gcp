package gcp

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
)

// GiBToBytes converting GB to bytes
const GiBToBytes = 1024 * 1024 * 1024

// TiBToGiB converting TiB to GiB
const TiBToGiB = 1024

func resourceGCPVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPVolumeCreate,
		Read:   resourceGCPVolumeRead,
		Delete: resourceGCPVolumeDelete,
		Update: resourceGCPVolumeUpdate,
		Exists: resourceGCPVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceVolumeCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_dp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"service_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "premium", "extreme"}, true),
			},
			"volume_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shared_vpc_project_number": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]+$"), "shared_vpc_project_number must be a numerical project number"),
			},
			"mount_points": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"snapshot_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"daily_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						"hourly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						"monthly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_of_month": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "1",
									},
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						"weekly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "Sunday",
									},
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
					},
				},
			},
			"export_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"allowed_clients": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"has_root_access": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "true",
										ValidateFunc: validation.StringInSlice([]string{"true", "false", "on", "off"}, true),
									},
									"kerberos5_readonly": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"kerberos5_readwrite": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"kerberos5i_readonly": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"kerberos5i_readwrite": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"kerberos5p_readonly": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"kerberos5p_readwrite": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"nfsv3": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"checked": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"nfsv4": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"checked": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"delete_on_creation_error": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"software", "hardware"}, true),
			},
			"regional_ha": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"snap_reserve": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"snapshot_directory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"smb_share_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"encrypt_data", "browsable", "changenotify", "non_browsable", "oplocks", "showsnapshot", "show_previous_versions", "continuously_available", "access_based_enumeration"}, true),
				},
			},
			"unix_permissions": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// TranslateServiceLevelState2API to translate service level state based on the setup value due to the API bugs
// resource value: API call value
// standard      : low
// premium       : medium
// extreme       : extreme
func TranslateServiceLevelState2API(slevel string) string {
	var apiValue = slevel
	if slevel == "standard" {
		apiValue = "low"
	} else if slevel == "premium" {
		apiValue = "medium"
	}

	return apiValue
}

func resourceGCPVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %v", d.Get("name").(string))

	client := meta.(*Client)

	volume := volumeRequest{}

	volume.Name = d.Get("name").(string)
	volume.Region = d.Get("region").(string)
	volume.Network = d.Get("network").(string)
	protocols := d.Get("protocol_types")
	for _, protocol := range protocols.([]interface{}) {
		if protocol.(string) == "SMB" {
			volume.ProtocolTypes = append(volume.ProtocolTypes, "CIFS")
		} else {
			volume.ProtocolTypes = append(volume.ProtocolTypes, protocol.(string))
		}
	}
	// size in 1 GiB increments, api takes in bytes only
	volume.Size = d.Get("size").(int) * GiBToBytes

	if v, ok := d.GetOk("service_level"); ok {
		slevel := v.(string)
		volume.ServiceLevel = TranslateServiceLevelState2API(slevel)
	}

	if v, ok := d.GetOk("snapshot_policy"); ok {
		if len(v.([]interface{})) > 0 {
			policy := v.([]interface{})[0].(map[string]interface{})
			volume.SnapshotPolicy = expandSnapshotPolicy(policy)
		}
	}

	if v, ok := d.GetOk("volume_path"); ok {
		volume.CreationToken = v.(string)
	}

	var volType string
	dpType := d.Get("type_dp").(bool)

	if dpType == true {
		volType = "DataProtectionVolumes"
	} else {
		volType = "Volumes"
	}

	if v, ok := d.GetOk("export_policy"); ok {
		policy := v.(*schema.Set)
		if policy.Len() > 0 {
			volume.ExportPolicy = expandExportPolicy(policy)
		}
	}

	if v, ok := d.GetOk("shared_vpc_project_number"); ok {
		volume.SharedVpcProjectNumber = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		volume.Zone = v.(string)
	}

	if v, ok := d.GetOk("storage_class"); ok {
		volume.StorageClass = v.(string)
	}

	if v, ok := d.GetOk("regional_ha"); ok {
		volume.RegionalHA = v.(bool)
	}

	// If storage class is 'software', zone or regionalHA is mandatory
	if volume.StorageClass == "software" && ((volume.Zone == "" && volume.RegionalHA == false) || (volume.Zone != "" && volume.RegionalHA == true)) {
		log.Print("Error creating volume")
		return fmt.Errorf("If storage_class is software, zone or RegionalHA is mandatory")
	}

	if v, ok := d.GetOk("snap_reserve"); ok {
		volume.SnapReserve = v.(int)
	}

	if v, ok := d.GetOk("unix_permissions"); ok {
		volume.UnixPermissions = v.(string)
	}

	volume.SnapshotDirectory = d.Get("snapshot_directory").(bool)

	if v, ok := d.GetOk("smb_share_settings"); ok {
		for _, setting := range v.(*schema.Set).List() {
			volume.SmbShareSettings = append(volume.SmbShareSettings, setting.(string))
		}
	}

	var res createVolumeResult
	var err error
	res, err = client.createVolume(&volume, volType)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}

	var volumeRes volumeResult
	time.Sleep(5 * time.Second)
	volume.Network = d.Get("network").(string)
	volumeRes, err = validateVolumeExistsAfterCreate(client, volume, res.Name.JobID.VolID, volType)
	if err != nil {
		return err
	}
	d.SetId(volumeRes.VolumeID)
	if volumeRes.LifeCycleState == "available" {
		return resourceGCPVolumeRead(d, meta)
	}
	volumeRes, err = waitForVolumeCreationComplete(client, volumeRes)
	if err != nil {
		return err
	}
	// if volume's state is error, delete the volume and retry for twice. If the operation still fails, return error.
	if volumeRes.LifeCycleState == "error" {
		retries := 2
		for retries > 0 && volumeRes.LifeCycleState == "error" {
			deleteErr := resourceGCPVolumeDelete(d, meta)
			if deleteErr != nil {
				return fmt.Errorf("failed to delete volume in error state after creation. %s", deleteErr.Error())
			}
			volume.Network = d.Get("network").(string)
			res, err = client.createVolume(&volume, volType)
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
			volume.Network = d.Get("network").(string)
			volumeRes, err = validateVolumeExistsAfterCreate(client, volume, res.Name.JobID.VolID, volType)
			if err != nil {
				return err
			}
			d.SetId(volumeRes.VolumeID)
			volumeRes, err = waitForVolumeCreationComplete(client, volumeRes)
			if err != nil {
				return err
			}
			if volumeRes.LifeCycleState == "available" {
				return resourceGCPVolumeRead(d, meta)
			}
			timeSleep := time.Duration(nextRandomInt(5, 10)) * time.Second
			time.Sleep(timeSleep)
			retries--
		}
		if d.Get("delete_on_creation_error").(bool) {
			deleteErr := resourceGCPVolumeDelete(d, meta)
			if deleteErr != nil {
				return fmt.Errorf("failed to delete volume in error state after creation. %s", deleteErr.Error())
			}
			return fmt.Errorf("%v. Volume in error state is deleted", volumeRes.LifeCycleStateDetails)
		}
		return fmt.Errorf("%v", volumeRes.LifeCycleStateDetails)
	}
	return resourceGCPVolumeRead(d, meta)
}

// Wait up to 15 minutes for volume creation to complete.
func waitForVolumeCreationComplete(client *Client, volumeRes volumeResult) (volumeResult, error) {
	waitSeconds := 900    // first volume creation can take 11 minutes
	threshold := 900 - 60 // when to warn
	elapsed := time.Duration(0)
	var err error
	for waitSeconds > 0 && volumeRes.LifeCycleState == "creating" {
		timeSleep := time.Duration(nextRandomInt(20, 30))
		time.Sleep(timeSleep * time.Second)
		elapsed = elapsed + timeSleep
		volumeRes, err = client.getVolumeByID(volumeRequest{Region: volumeRes.Region, VolumeID: volumeRes.VolumeID})
		if err != nil {
			return volumeResult{}, err
		}
		if waitSeconds < threshold {
			threshold = threshold - 60
			log.Printf("Volume creation still in progress after %d seconds.\n", elapsed)
		}
		waitSeconds = waitSeconds - int(timeSleep)
	}
	return volumeRes, nil
}

// A bug might be presented in the API. A volume creation request is acknowledged(volume ID is returned), but get volume by ID doesn't find any result.
// A temporary fix is to send the create request again.
func validateVolumeExistsAfterCreate(client *Client, volume volumeRequest, volumeID string, volType string) (volumeResult, error) {
	volumeRes, err := client.getVolumeByID(volumeRequest{Region: volume.Region, VolumeID: volumeID})
	var res createVolumeResult
	network := volume.Network
	retries := 3
	if err != nil {
		for err != nil && err.Error() == "code: 404, message: Error describing volume - Volume not found" && retries > 0 {
			time.Sleep(20 * time.Second)
			volume.Network = network
			res, err = client.createVolume(&volume, volType)
			if err != nil {
				return volumeResult{}, err
			}
			volumeRes, err = client.getVolumeByID(volumeRequest{Region: volume.Region, VolumeID: res.Name.JobID.VolID})
			retries--
		}
		if err != nil {
			return volumeResult{}, err
		}
	}
	return volumeRes, nil
}

func resourceGCPVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume: %#v", d)
	client := meta.(*Client)

	volume := volumeRequest{}

	volume.Region = d.Get("region").(string)

	id := d.Id()
	volume.VolumeID = id

	var res volumeResult
	res, err := client.getVolumeByID(volume)
	if err != nil {
		return err
	}

	waitSeconds := 300
	for waitSeconds > 0 && (res.LifeCycleState == "creating" || res.LifeCycleState == "deleting" || res.LifeCycleState == "updating") {
		time.Sleep(20)
		res, err = client.getVolumeByID(volumeRequest{Region: volume.Region, VolumeID: id})
		if err != nil {
			return err
		}
		waitSeconds = waitSeconds - 10
	}

	if res.VolumeID != id {
		return fmt.Errorf("Expected Volume ID %v, Response contained Volume ID %v", id, res.VolumeID)
	}

	if res.LifeCycleState == "error" {
		return fmt.Errorf("Volume with name: %v and id: %v is in error state. Please manually delete the volume, make sure the config is correct and run terraform apply again. LifeCycleStateDetails: %v",
			res.Name, res.VolumeID, res.LifeCycleStateDetails)
	} else if res.LifeCycleState == "disabled" {
		return fmt.Errorf("Volume with name: %v and id: %v is in disabled state. Please manually enable the volume and runn terraform apply again. LifeCycleStateDetails: %v",
			res.Name, res.VolumeID, res.LifeCycleStateDetails)
	} else if res.LifeCycleState == "deleted" {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", res.Name); err != nil {
		return fmt.Errorf("Error reading volume name: %s", err)
	}

	if err := d.Set("size", res.Size/GiBToBytes); err != nil {
		return fmt.Errorf("Error reading volume size: %s", err)
	}

	log.Printf("**** API response service level is %s", res.ServiceLevel)
	/*
		There is a bug on API. Translate the response with right value
		API Respose:  Real Value
		basic      :  standard
		standard   :  premium
		extreme    :  extreme
	*/
	var slevel = res.ServiceLevel

	if res.ServiceLevel == "basic" {
		slevel = "standard"
	} else if res.ServiceLevel == "standard" {
		slevel = "premium"
	}

	if err := d.Set("service_level", slevel); err != nil {
		return fmt.Errorf("Error reading volume service_level: %s", err)
	}
	for i, protocol := range res.ProtocolTypes {
		if protocol == "CIFS" {
			res.ProtocolTypes[i] = "SMB"
		}
	}
	if err := d.Set("protocol_types", res.ProtocolTypes); err != nil {
		return fmt.Errorf("Error reading volume protocol_types: %s", err)
	}
	if err := d.Set("volume_path", res.CreationToken); err != nil {
		return fmt.Errorf("Error reading volume path or Creation Token: %s", err)
	}
	// res.Network either contains simple network name or
	// projects/${HOST_PROJECT_ID}/global/networks/${SHARED_VPC_NAME}, usually (but not exclusively) for shared VPC
	nws := strings.Split(res.Network, "/")
	var network string
	if len(nws) == 1 {
		// standalone project network
		network = nws[0]
	} else if len(nws) == 5 {
		// long network path
		network = nws[4]
		// if network path contains different projectId than our project, it is shared-VPC
		if nws[1] != client.Project {
			if err := d.Set("shared_vpc_project_number", nws[1]); err != nil {
				return fmt.Errorf("Error reading shared_vpc_project_number: %s", err)
			}
		}
	} else {
		return fmt.Errorf("Error returned network path invalid: %s", res.Network)
	}
	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("Error reading volume network: %s", err)
	}
	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("Error reading volume region: %s", err)
	}
	snapshotPolicy := flattenSnapshotPolicy(res.SnapshotPolicy)
	exportPolicy := flattenExportPolicy(res.ExportPolicy)
	if err := d.Set("snapshot_policy", snapshotPolicy); err != nil {
		return fmt.Errorf("Error reading volume snapshot_policy: %s", err)
	}
	if len(res.ExportPolicy.Rules) > 0 {
		if err := d.Set("export_policy", exportPolicy); err != nil {
			return fmt.Errorf("Error reading volume export_policy: %s", err)
		}
	} else {
		a := schema.NewSet(schema.HashString, []interface{}{})
		if err := d.Set("export_policy", a); err != nil {
			return fmt.Errorf("Error reading volume export_policy: %s", err)
		}
	}
	mountPoints := flattenMountPoints(res.MountPoints)
	if err := d.Set("mount_points", mountPoints); err != nil {
		return fmt.Errorf("Error reading volume mount_points: %s", err)
	}
	if _, ok := d.GetOk("zone"); ok {
		if err := d.Set("zone", res.Zone); err != nil {
			return fmt.Errorf("Error reading volume zone: %s", err)
		}
	}
	if _, ok := d.GetOk("storage_class"); ok {
		if err := d.Set("storage_class", res.StorageClass); err != nil {
			return fmt.Errorf("Error reading volume storage_class: %s", err)
		}
	}
	if err := d.Set("snap_reserve", res.SnapReserve); err != nil {
		return fmt.Errorf("Error reading volume snap_reserve: %s", err)
	}
	if err := d.Set("snapshot_directory", res.SnapshotDirectory); err != nil {
		return fmt.Errorf("Error reading volume snapshot_directory: %s", err)
	}
	if v, ok := d.GetOk("smb_share_settings"); ok {
		// There are a few default values in API, which means the API sets these values even they aren't specified in creation.
		// The default values: "oplocks", "changenotify", "showsnapshot", "show_previous_versions", "browsable".
		// Also note: If "continuously_available" is specified and "changenotify" is not, changenotify won't be set.
		// but it doesn't mean they are mutually exclusive. "changenotify" can still be specified.
		// "browserable" and "non_browserable" are mutually exclusive.
		// It's reasonable that Users only care about the smb_share_settings they specifed. The not specified smb_share_settings
		// are ignored no matter enabled or not.
		// Compare the smb_share_settings in local config and the API response, then set the intersection of the two lists.
		currentSmbSettings := make([]string, 0)
		for _, localSmbSetting := range v.(*schema.Set).List() {
			for _, apiSmbSetting := range res.SmbShareSettings {
				if localSmbSetting.(string) == apiSmbSetting {
					currentSmbSettings = append(currentSmbSettings, apiSmbSetting)
				}
			}
		}
		if err := d.Set("smb_share_settings", currentSmbSettings); err != nil {
			return fmt.Errorf("Error reading volume smb_share_settings: %s", err)
		}
	}
	if _, ok := d.GetOk("unix_permissions"); ok {
		if err := d.Set("unix_permissions", res.UnixPermissions); err != nil {
			return fmt.Errorf("Error reading volume unix_permissions: %S", err)
		}
	}
	return nil
}

func resourceGCPVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume: %#v", d)

	volume := volumeRequest{}

	volume.Region = d.Get("region").(string)
	client := meta.(*Client)

	id := d.Id()
	volume.VolumeID = id

	deleteErr := client.deleteVolume(volume)
	if deleteErr != nil {
		return deleteErr
	}

	getVolume, err := client.getVolumeByID(volume)
	if err != nil {
		return err
	}
	if getVolume.LifeCycleState == "deleted" {
		return nil
	} else if getVolume.LifeCycleState == "deleting" {
		waitTime := 300
		for waitTime > 0 {
			time.Sleep(20 * time.Second)
			waitTime = waitTime - 20
			getVolume, err = client.getVolumeByID(volume)
			if err != nil {
				return err
			}
			if getVolume.LifeCycleState == "deleted" {
				return nil
			}
			if getVolume.LifeCycleState == "error" {
				break
			}
		}
		// if volume is in error state when deleting, retry.
	}
	if getVolume.LifeCycleState == "error" {
		retries := 3
		for getVolume.LifeCycleState == "error" && retries > 0 {
			time.Sleep(time.Duration(nextRandomInt(5, 20)) * time.Second)
			deleteErr := client.deleteVolume(volume)
			if deleteErr != nil {
				return deleteErr
			}
			getVolume, err = client.getVolumeByID(volume)
			if err != nil {
				return err
			}
		}
		if getVolume.LifeCycleState == "error" {
			return fmt.Errorf("error deleting volume with id: %s, name: %s; %s", getVolume.VolumeID, getVolume.Name, getVolume.LifeCycleStateDetails)
		}
	}

	return nil
}

func resourceGCPVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume: %#v", d)
	client := meta.(*Client)

	volume := volumeRequest{}

	id := d.Id()
	volume.VolumeID = id
	volume.Region = d.Get("region").(string)
	var res volumeResult
	res, err := client.getVolumeByID(volume)
	if err != nil {
		if err, ok := err.(*restapi.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
			return false, err
		}
		return false, err
	}

	if res.VolumeID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}

func resourceGCPVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume: %#v\n", d)
	makechange := 0
	client := meta.(*Client)
	volume := volumeRequest{}
	volume.VolumeID = d.Id()
	volume.Region = d.Get("region").(string)
	volume.Name = d.Get("name").(string)
	// size is always required.
	volume.Size = d.Get("size").(int) * GiBToBytes
	volume.SnapReserve = d.Get("snap_reserve").(int)
	volume.SnapshotDirectory = d.Get("snapshot_directory").(bool)

	if d.HasChange("name") {
		makechange = 1
	}

	if d.HasChange("size") {
		makechange = 1
	}

	if d.HasChange("snap_reserve") {
		makechange = 1
	}

	if d.HasChange("snapshot_directory") {
		makechange = 1
	}

	if d.HasChange("snapshot_policy") {
		if len(d.Get("snapshot_policy").([]interface{})) > 0 {
			policy := d.Get("snapshot_policy").([]interface{})[0].(map[string]interface{})
			volume.SnapshotPolicy = expandSnapshotPolicy(policy)
			makechange = 1
		}
	}

	if d.HasChange("export_policy") {
		policy := d.Get("export_policy").(*schema.Set)
		volume.ExportPolicy = expandExportPolicy(policy)
		makechange = 1
	}

	if d.HasChange("service_level") {
		o, n := d.GetChange("service_level")
		slevel := n.(string)
		oslevel := o.(string)

		log.Printf("Updating volume: service_level old=%v new=%v\n", oslevel, slevel)
		volume.ServiceLevel = TranslateServiceLevelState2API(slevel)
		makechange = 1
	}

	if d.HasChange("smb_share_settings") {
		if v, ok := d.GetOk("smb_share_settings"); ok {
			for _, setting := range v.(*schema.Set).List() {
				volume.SmbShareSettings = append(volume.SmbShareSettings, setting.(string))
			}
		}
		makechange = 1
	}

	if d.HasChange("unix_permissions") {
		makechange = 1
		volume.UnixPermissions = d.Get("unix_permissions").(string)
	}

	if makechange == 1 {
		log.Println("Make change on volume")
		err := client.updateVolume(volume)
		if err != nil {
			return err
		}
	} else {
		log.Println("NOT updateVolume")
	}

	return resourceGCPVolumeRead(d, meta)
}

func resourceVolumeCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("storage_class") {
		current, expect := diff.GetChange("storage_class")
		if current.(string) == "" {
			if expect.(string) == "software" {
				if diff.Get("service_level").(string) != "standard" {
					return fmt.Errorf("service_level must be standard when storage_class is software")
				}
			} else if expect.(string) == "hardware" {
				if v, ok := diff.GetOk("regional_ha"); ok {
					if v.(bool) == true {
						return fmt.Errorf("regional_ha is not supported when storage_class is hardware")
					}
				}
			}
		}
	}
	return nil
}
