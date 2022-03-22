package gcp

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGCPKMSConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPKMSConfigCreate,
		Read:   resourceGCPKMSConfigRead,
		Delete: resourceGCPKMSConfigDelete,
		Update: resourceGCPKMSConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			//these available fields are required for create and update.
			"key_ring_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_ring_location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGCPKMSConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	kms := kmsConfig{}
	kms.KeyRing = d.Get("key_ring_name").(string)
	kms.KeyName = d.Get("key_name").(string)
	kms.KeyRingLocation = d.Get("key_ring_location").(string)
	kms.Network = d.Get("network").(string)

	if v, ok := d.GetOk("key_project_id"); ok {
		kms.KeyProjectID = v.(string)
	}

	res, err := client.createKMSConfig(&kms)
	if err != nil {
		log.Print("Error creating kms config")
		return err
	}
	d.SetId(res.ID)

	log.Printf("Created KMS in region: %v", kms.KeyRingLocation)

	return resourceGCPKMSConfigRead(d, meta)
}

func resourceGCPKMSConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	id := d.Id()
	kmsConfig := kmsConfig{}
	kmsConfig.KeyRingLocation = d.Get("key_ring_location").(string)
	kmsConfig.ID = d.Id()
	res, err := client.getKMSConfig(&kmsConfig)
	if err != nil {
		return err
	}
	if res.ID != id {
		return fmt.Errorf("Expected kms with id: %v, Response contained kms with id: %v",
			d.Id(), res.ID)
	}

	if err := d.Set("key_name", res.KeyName); err != nil {
		return fmt.Errorf("Error reading key name: %s", err)
	}

	if err := d.Set("key_ring_name", res.KeyRing); err != nil {
		return fmt.Errorf("Error reading key ring name: %s", err)
	}

	if err := d.Set("key_ring_location", res.KeyRingLocation); err != nil {
		return fmt.Errorf("Error reading key ring location: %s", err)
	}

	if _, ok := d.GetOk("key_projet_id"); ok {
		if err := d.Set("key_project_id", res.KeyProjectID); err != nil {
			return fmt.Errorf("Error reading key project id: %s", err)
		}
	}

	if err := d.Set("network", res.Network); err != nil {
		return fmt.Errorf("Error reading network id: %s", err)
	}

	return nil
}

func resourceGCPKMSConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting kms: %#v", d)
	client := meta.(*Client)
	kms := kmsConfig{}
	kms.KeyRingLocation = d.Get("key_ring_location").(string)
	kms.ID = d.Id()
	_, deleteErr := client.deleteKMSConfig(&kms)
	if deleteErr != nil {
		return deleteErr
	}
	d.SetId("")

	return nil
}

func resourceGCPKMSConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("updating kms: %#v", d)
	client := meta.(*Client)
	kms := kmsConfig{}
	// all of the following are required for API: update.
	kms.KeyName = d.Get("key_name").(string)
	kms.KeyRing = d.Get("key_ring_name").(string)
	kms.KeyRingLocation = d.Get("key_ring_location").(string)
	kms.Network = d.Get("network").(string)
	kms.ID = d.Id()

	if v, ok := d.GetOk("key_project_id"); ok {
		kms.KeyProjectID = v.(string)
	}

	_, err := client.updateKMSConfig(&kms)
	if err != nil {
		return err
	}
	return resourceGCPKMSConfigRead(d, meta)
}
