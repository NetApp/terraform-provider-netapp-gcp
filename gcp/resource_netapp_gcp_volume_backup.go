package gcp

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
)

func resourceGCPVolumeBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPVolumeBackupCreate,
		Read:   resourceGCPVolumeBackupRead,
		Delete: resourceGCPVolumeBackupDelete,
		Exists: resourceGCPVolumeBackupExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceGCPVolumeBackupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume backup: %#v", d)

	client := meta.(*Client)

	volumeBackup := createVolumeBackupRequest{}

	volumeBackup.Name = d.Get("name").(string)
	volumeBackup.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = volumeBackup.Region

	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	// Check the volume status. Start creating backup when volume is ready to use
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
			volumeBackup.VolumeID = volresult.VolumeID
			break
		}
	}

	res, err := client.createVolumeBackup(&volumeBackup)
	if err != nil {
		log.Print("Error creating VolumeBackup")
		return err
	}

	d.SetId(res.Name.JobID.VolumeBackupID)
	log.Printf("Created VolumeBackup: %v", volumeBackup.Name)

	return resourceGCPVolumeBackupRead(d, meta)
}

func resourceGCPVolumeBackupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading VolumeBackup: %#v", d)
	client := meta.(*Client)

	volumeBackup := listVolumeBackupRequest{}

	volumeBackup.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = volumeBackup.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return err
	}

	volumeBackup.VolumeID = volresult.VolumeID

	id := d.Id()
	volumeBackup.VolumeBackupID = id
	var res listVolumeBackupResult
	res, err = client.getVolumeBackupByID(volumeBackup)
	if err != nil {
		log.Print("Error getting VolumeBackup")
		return err
	}

	if res.VolumeBackupID != id {
		return fmt.Errorf("Expected VolumeBackup ID %v, Response contained VolumeBackup ID %v", id, res.VolumeBackupID)
	}

	return nil
}

func resourceGCPVolumeBackupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting VolumeBackup: %#v", d)

	client := meta.(*Client)

	volumeBackup := deleteVolumeBackupRequest{}

	volumeBackup.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = volumeBackup.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return err
	}

	volumeBackup.VolumeID = volresult.VolumeID

	id := d.Id()
	volumeBackup.VolumeBackupID = id

	deleteErr := client.deleteVolumeBackup(volumeBackup)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceGCPVolumeBackupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of VolumeBackup: %#v", d)
	client := meta.(*Client)

	volumeBackup := listVolumeBackupRequest{}

	id := d.Id()
	volumeBackup.VolumeBackupID = id
	volumeBackup.Region = d.Get("region").(string)

	volume := volumeRequest{}
	volume.Region = volumeBackup.Region
	volume.Name = d.Get("volume_name").(string)
	volume.CreationToken = d.Get("creation_token").(string)

	volresult, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		log.Print("Error getting volume ID")
		return false, err
	}

	volumeBackup.VolumeID = volresult.VolumeID

	var res listVolumeBackupResult
	res, err = client.getVolumeBackupByID(volumeBackup)
	if err != nil {
		if err, ok := err.(*restapi.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if res.VolumeBackupID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
