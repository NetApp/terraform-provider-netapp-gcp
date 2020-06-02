package gcp

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
)

func resourceGCPSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPSnapshotCreate,
		Read:   resourceGCPSnapshotRead,
		Delete: resourceGCPSnapshotDelete,
		Exists: resourceGCPSnapshotExists,
		Update: resourceGCPSnapshotUpdate,
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
				ForceNew: true,
			},
			"volume_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"creation_token": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGCPSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating snapshot: %#v", d)

	client := meta.(*Client)

	snapshot := createSnapshotRequest{}

	snapshot.Name = d.Get("name").(string)
	snapshot.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = snapshot.Region

	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	// Check the volume status. Start creating snapshot when volume is ready to use
	retries := 0
	for {
		volresult, err := client.getVolumeByNameOrCreationToken(volume)
		if err != nil {
			log.Print("Error getting volume ID")
			return err
		}
		if volresult.LifeCycleStateDetails != "Available for use" {
			if retries < 3 {
				log.Printf("Volume %s is not ready. Wait for 5 seconds and check again.\n", volume.Name)
				time.Sleep(5 * time.Second)
				retries++
			} else {
				log.Printf("Volume %s is not ready.\n", volume.Name)
				return err
			}
		} else {
			snapshot.VolumeID = volresult.VolumeID
			break
		}
	}

	res, err := client.createSnapshot(&snapshot)
	if err != nil {
		log.Print("Error creating snapshot")
		return err
	}

	d.SetId(res.Name.JobID.SnapshotID)
	log.Printf("Created snapshot: %v", snapshot.Name)

	return resourceGCPSnapshotRead(d, meta)
}

func resourceGCPSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading snapshot: %#v", d)
	client := meta.(*Client)

	snapshot := listSnapshotRequest{}

	snapshot.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = snapshot.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return err
	}

	snapshot.VolumeID = volresult.VolumeID

	id := d.Id()
	snapshot.SnapshotID = id
	var res listSnapshotResult
	res, err = client.getSnapshotByID(snapshot)
	if err != nil {
		log.Print("Error getting Snapshot")
		return err
	}

	if res.SnapshotID != id {
		return fmt.Errorf("Expected Snapshot ID %v, Response contained Snapshot ID %v", id, res.SnapshotID)
	}

	return nil
}

func resourceGCPSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting snapshot: %#v", d)

	client := meta.(*Client)

	snapshot := deleteSnapshotRequest{}

	snapshot.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = snapshot.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return err
	}

	snapshot.VolumeID = volresult.VolumeID

	id := d.Id()
	snapshot.SnapshotID = id

	deleteErr := client.deleteSnapshot(snapshot)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceGCPSnapshotExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of snapshot: %#v", d)
	client := meta.(*Client)

	snapshot := listSnapshotRequest{}

	id := d.Id()
	snapshot.SnapshotID = id
	snapshot.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = snapshot.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return false, err
	}

	snapshot.VolumeID = volresult.VolumeID

	var res listSnapshotResult
	res, err = client.getSnapshotByID(snapshot)
	if err != nil {
		if err, ok := err.(*restapi.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if res.SnapshotID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}

func resourceGCPSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating snapshot: %#v", d)

	client := meta.(*Client)

	snapshot := updateSnapshotRequest{}
	id := d.Id()
	snapshot.SnapshotID = id
	snapshot.Name = d.Get("name").(string)
	snapshot.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = snapshot.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return err
	}

	snapshot.VolumeID = volresult.VolumeID

	err = client.updateSnapshot(snapshot)
	if err != nil {
		return err
	}

	log.Printf("Updated snapshot: %v", snapshot.Name)

	return resourceGCPSnapshotRead(d, meta)
}
