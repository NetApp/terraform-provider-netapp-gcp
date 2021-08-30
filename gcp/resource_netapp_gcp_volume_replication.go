package gcp

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
)

func resourceGCPVolumeReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPVolumeReplicationCreate,
		Read:   resourceGCPVolumeReplicationRead,
		Delete: resourceGCPVolumeReplicationDelete,
		Update: resourceGCPVolumeReplicationUpdate,
		Exists: resourceGCPVolumeReplicationExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"destination_volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"MirrorAllSnapshots", "MirrorLatest", "MirrorAndVault"}, true),
			},
			"schedule": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"10minutely", "hourly", "daily", "weekly", "monthly"}, true),
			},
			"bandwidth": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"64Mbps", "128Mbps", "256Mbps"}, true),
			},
		},
	}
}

func resourceGCPVolumeReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume replication: %#v", d)

	client := meta.(*Client)

	replica := volumeReplicationRequest{}

	replica.Name = d.Get("name").(string)
	replica.Region = d.Get("region").(string)

	if v, ok := d.GetOk("destination_volume_id"); ok {
		replica.DestinationVolumeID = v.(string)
	}
	if v, ok := d.GetOk("source_volume_id"); ok {
		replica.SourceVolumeID = v.(string)
	}
	if v, ok := d.GetOk("remote_region"); ok {
		replica.RemoteRegion = v.(string)
	}
	if v, ok := d.GetOk("policy"); ok {
		replica.Policy = v.(string)
	}
	if v, ok := d.GetOk("schedule"); ok {
		replica.Schedule = v.(string)
	}
	if v, ok := d.GetOk("mirror_state"); ok {
		replica.MirrorState = v.(string)
	}
	if v, ok := d.GetOk("endpoint_type"); ok {
		replica.EndpointType = v.(string)
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		replica.Bandwidth = v.(string)
	}
	log.Printf("here here here here here: %#v", replica)

	res, err := client.createVolumeReplication(&replica)
	if err != nil {
		log.Print("Error creating volume replication")
		return err
	}

	d.SetId(res.ReplicationID)
	log.Printf("Created volume replication: %v", res.Name)

	return resourceGCPVolumeReplicationRead(d, meta)
}

func resourceGCPVolumeReplicationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume replication: %#v", d)
	client := meta.(*Client)

	replication := volumeReplicationRequest{}

	replication.Region = d.Get("region").(string)

	id := d.Id()
	replication.ReplicationID = id

	var res volumeReplicationResult
	for {
		var replica volumeReplicationResult
		replica, err := client.getVolumeReplicationByID(replication)
		if err != nil {
			return err
		}

		if replica.ReplicationID != id {
			return fmt.Errorf("Expected replication ID %v, Response contained replication ID %v", id, res.ReplicationID)
		}

		if replica.LifeCycleState == "error" {
			return fmt.Errorf("Volume replication %v is in %v state. Please check the setup. Will delete the volume replication",
				replica.ReplicationID, replica.LifeCycleState)
		} else if replica.LifeCycleState == "available" {
			res = replica
			break
		} else {
			time.Sleep(time.Duration(2) * time.Second)
		}
	}

	if err := d.Set("destination_volume_id", res.DestinationVolumeID); err != nil {
		return fmt.Errorf("Error reading destination volume id: %s", err)
	}

	if err := d.Set("source_volume_id", res.SourceVolumeID); err != nil {
		return fmt.Errorf("Error reading source volume id: %s", err)
	}

	if err := d.Set("remote_region", res.RemoteRegion); err != nil {
		return fmt.Errorf("Error reading source volume id: %s", err)
	}

	if err := d.Set("endpoint_type", res.EndpointType); err != nil {
		return fmt.Errorf("Error reading source volume endpoint_type: %s", err)
	}

	if err := d.Set("policy", res.Policy); err != nil {
		return fmt.Errorf("Error reading policy: %s", err)
	}

	if err := d.Set("schedule", res.Schedule); err != nil {
		return fmt.Errorf("Error reading schedule: %s", err)
	}

	if err := d.Set("bandwidth", res.Bandwidth); err != nil {
		return fmt.Errorf("Error reading bandwidth: %s", err)
	}
	return nil
}

func resourceGCPVolumeReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume replication: %#v", d)

	replica := volumeReplicationRequest{}

	replica.Region = d.Get("region").(string)
	client := meta.(*Client)

	id := d.Id()
	replica.ReplicationID = id

	err := client.breakVolumeReplication(&replica)
	if err != nil {
		return err
	}
	err = client.deleteVolumeReplication(&replica)
	if err != nil {
		return err
	}

	return nil
}

func resourceGCPVolumeReplicationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume replication: %#v", d)
	client := meta.(*Client)

	replica := volumeReplicationRequest{}

	id := d.Id()
	replica.ReplicationID = id
	replica.Region = d.Get("region").(string)
	var res volumeReplicationResult
	res, err := client.getVolumeReplicationByID(replica)
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

	if res.ReplicationID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}

func resourceGCPVolumeReplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume replication: %#v\n", d)
	client := meta.(*Client)
	replica := volumeReplicationRequest{}
	replica.ReplicationID = d.Id()

	if d.HasChange("schedule") {
		replica.Schedule = d.Get("schedule").(string)
	}

	if d.HasChange("name") {
		replica.Name = d.Get("name").(string)
	}

	if d.HasChange("policy") {
		replica.Policy = d.Get("policy").(string)
	}

	if d.HasChange("bandwidth") {
		replica.Bandwidth = d.Get("bandwidth").(string)
	}

	err := client.updateVolumeReplication(&replica)
	if err != nil {
		return err
	}
	return resourceGCPVolumeRead(d, meta)
}
