package gcp

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
)

// GiBTobytes converting GB to bytes
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
				ValidateFunc: validation.StringInSlice([]string{"low", "medium", "high"}, true),
			},
			"snapshot_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"daily_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"hourly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"monthly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_of_month": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"weekly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"export_policy": {
				Type:     schema.TypeList,
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
									"nfsv3": {
										Type:     schema.TypeList,
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
										Type:     schema.TypeList,
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
		},
	}
}

func resourceGCPVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %#v", d)

	client := meta.(*Client)

	volume := createVolumeRequest{}

	volume.Name = d.Get("name").(string)
	volume.Region = d.Get("region").(string)
	volume.Network = d.Get("network").(string)
	protocols := d.Get("protocol_types")
	for _, protocol := range protocols.([]interface{}) {
		volume.ProtocolTypes = append(volume.ProtocolTypes, protocol.(string))
	}
	// size in 1 GiB increments, api takes in bytes only
	volume.Size = d.Get("size").(int) * GiBToBytes

	if v, ok := d.GetOk("service_level"); ok {
		volume.ServiceLevel = v.(string)
	}

	if v, ok := d.GetOk("snapshot_policy"); ok {
		if len(v.([]interface{})) > 0 {
			policy := v.([]interface{})[0].(map[string]interface{})
			volume.SnapshotPolicy = expandSnapshotPolicy(policy)
		}
	}

	if v, ok := d.GetOk("export_policy"); ok {
		if len(v.([]interface{})) > 0 {
			policy := v.([]interface{})[0].(map[string]interface{})
			volume.ExportPolicy = expandExportPolicy(policy)
		}
	}

	res, err := client.createVolume(&volume)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}

	d.SetId(res.Name.JobID.VolID)
	log.Printf("Created volume: %v", volume.Name)

	return resourceGCPVolumeRead(d, meta)
}

func resourceGCPVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume: %#v", d)
	client := meta.(*Client)

	volume := listVolumesRequest{}

	volume.Region = d.Get("region").(string)

	id := d.Id()
	volume.VolumeID = id
	var res listVolumeResult
	res, err := client.getVolumeByID(volume)
	if err != nil {
		return err
	}

	if res.VolumeID != id {
		return fmt.Errorf("Expected VOlume ID %v, Response contained Volume ID %v", id, res.VolumeID)
	}

	return nil
}

func resourceGCPVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume: %#v", d)

	volume := deleteVolumeRequest{}

	volume.Region = d.Get("region").(string)
	client := meta.(*Client)

	id := d.Id()
	volume.VolumeID = id

	deleteErr := client.deleteVolume(volume)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceGCPVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume: %#v", d)
	client := meta.(*Client)

	volume := listVolumesRequest{}

	id := d.Id()
	volume.VolumeID = id
	volume.Region = d.Get("region").(string)
	var res listVolumeResult
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
	log.Printf("Updating volume: %#v", d)
	client := meta.(*Client)
	volume := updateVolumeRequest{}
	volume.VolumeID = d.Id()
	volume.Region = d.Get("region").(string)
	volume.Name = d.Get("name").(string)
	// size is always required.
	volume.Size = d.Get("size").(int) * GiBToBytes

	if d.HasChange("snapshot_policy") {
		if len(d.Get("snapshot_policy").([]interface{})) > 0 {
			policy := d.Get("snapshot_policy").([]interface{})[0].(map[string]interface{})
			volume.SnapshotPolicy = expandSnapshotPolicy(policy)
		}
	}

	if d.HasChange("export_policy") {
		if len(d.Get("export_policy").([]interface{})) > 0 {
			policy := d.Get("export_policy").([]interface{})[0].(map[string]interface{})
			volume.ExportPolicy = expandExportPolicy(policy)
		}
	}

	if d.HasChange("service_level") {
		volume.ServiceLevel = d.Get("service_level").(string)
	}

	err := client.updateVolume(volume)
	if err != nil {
		log.Print("updateVolume request failed")
		return err
	}

	return nil
}
