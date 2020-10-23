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
				Default:      "medium",
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
							Type:     schema.TypeSet,
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
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

/*
	There is a bug on service level API. Translate based on the setup value.
	resource value: API call value
	standard      : low
	premium       : medium
	extreme       : extreme
*/
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

	log.Print(volType)

	if v, ok := d.GetOk("export_policy"); ok {
		policy := v.(*schema.Set)
		if policy.Len() > 0 {
			volume.ExportPolicy = expandExportPolicy(policy)
		}
	}

	if v, ok := d.GetOk("shared_vpc_project_number"); ok {
		volume.Shared_vpc_project_number = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		volume.Zone = v.(string)
	}

	if v, ok := d.GetOk("storage_class"); ok {
		volume.StorageClass = v.(string)
	}
	var res createVolumeResult
	var err error
	res, err = client.createVolume(&volume, volType)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}

	d.SetId(res.Name.JobID.VolID)
	var volumeRes volumeResult
	time.Sleep(5 * time.Second)
	volumeRes, err = client.getVolumeByID(volumeRequest{Region: volume.Region, VolumeID: res.Name.JobID.VolID})
	if err != nil {
		return err
	}
	if volumeRes.LifeCycleState == "available" {
		return resourceGCPVolumeRead(d, meta)
	}
	if volumeRes.LifeCycleState == "creating" {
		// wait for 5 minutes for completing volume creation.
		waitSeconds := 300
		for waitSeconds > 0 && volumeRes.LifeCycleState == "creating" {
			time.Sleep(10)
			volumeRes, err = client.getVolumeByID(volumeRequest{Region: volume.Region, VolumeID: res.Name.JobID.VolID})
			if err != nil {
				return err
			}
			waitSeconds = waitSeconds - 10
		}
	} else if volumeRes.LifeCycleState == "error" {
		if d.Get("delete_on_creation_error").(bool) {
			deleteErr := resourceGCPVolumeDelete(d, meta)
			if deleteErr != nil {
				return fmt.Errorf("failed to delete volume in error state after creation. %s", deleteErr.Error())
			}
			return fmt.Errorf("%v. Volume in error state is deleted", err.Error())
		}
		return fmt.Errorf("%v", err.Error())
	}

	return resourceGCPVolumeRead(d, meta)
}

func resourceGCPVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume: %#v", d)
	client := meta.(*Client)

	volume := volumeRequest{}

	volume.Region = d.Get("region").(string)

	id := d.Id()
	volume.VolumeID = id

	var res volumeResult
	for {
		var vol volumeResult
		vol, err := client.getVolumeByID(volume)
		if err != nil {
			return err
		}

		if vol.VolumeID != id {
			return fmt.Errorf("Expected Volume ID %v, Response contained Volume ID %v", id, res.VolumeID)
		}

		if vol.LifeCycleState == "error" {
			return fmt.Errorf("Volume with name: %v and id: %v is in %v state. Please check the setup. LifeCycleStateDetails: %v",
				vol.Name, vol.VolumeID, vol.LifeCycleState, vol.LifeCycleStateDetails)
		} else if vol.LifeCycleState == "available" {
			res = vol
			break
		} else {
			time.Sleep(time.Duration(2) * time.Second)
		}
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
	network := res.Network
	index := strings.Index(network, "networks/")
	if index > -1 {
		network = network[index+len("networks/"):]
	}
	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("Error reading volume network: %s", err)
	}
	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("Error reading volume region: %s", err)
	}
	snapshot_policy := flattenSnapshotPolicy(res.SnapshotPolicy)
	export_policy := flattenExportPolicy(res.ExportPolicy)
	if err := d.Set("snapshot_policy", snapshot_policy); err != nil {
		return fmt.Errorf("Error reading volume snapshot_policy: %s", err)
	}
	if len(res.ExportPolicy.Rules) > 0 {
		if err := d.Set("export_policy", export_policy); err != nil {
			return fmt.Errorf("Error reading volume export_policy: %s", err)
		}
	} else {
		a := schema.NewSet(schema.HashString, []interface{}{})
		if err := d.Set("export_policy", a); err != nil {
			return fmt.Errorf("Error reading volume export_policy: %s", err)
		}
	}
	mount_points := flattenMountPoints(res.MountPoints)
	if err := d.Set("mount_points", mount_points); err != nil {
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
		}
		// if volume is in error state when deleting, retry.
	} else if getVolume.LifeCycleState == "error" {
		time.Sleep(time.Duration(nextRandomInt(5, 20)) * time.Second)
		deleteErr := client.deleteVolume(volume)
		if deleteErr != nil {
			return deleteErr
		}
		getVolume, err := client.getVolumeByID(volume)
		if err != nil {
			return err
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

	if d.HasChange("name") {
		makechange = 1
	}

	if d.HasChange("size") {
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
