package gcp

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceGCPStoragePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPStoragePoolCreate,
		Read:   resourceGCPStoragePoolRead,
		Delete: resourceGCPStoragePoolDelete,
		Update: resourceGCPStoragePoolUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_level": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ZoneRedundantStandardSW", "StandardSW"}, true),
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"regional_ha": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"hardware", "software"}, true),
			},
			"shared_vpc_project_number": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]+$"), "shared_vpc_project_number must be a numerical project number"),
			},
			"billing_label": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceGCPStoragePoolCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating storage pool: %#v", d.Get("name").(string))
	client := meta.(*Client)
	pool := storagePool{}
	// required attributes
	pool.Region = d.Get("region").(string)
	pool.Name = d.Get("name").(string)
	pool.ServiceLevel = d.Get("service_level").(string)
	pool.SizeInBytes = d.Get("size").(int) * GiBToBytes
	pool.Network = d.Get("network").(string)
	// optional attributes
	if v, ok := d.GetOk("regional_ha"); ok {
		pool.RegionalHA = v.(bool)
	}
	if v, ok := d.GetOk("zone"); ok {
		pool.Zone = v.(string)
	}
	if v, ok := d.GetOk("storage_class"); ok {
		pool.StorageClass = v.(string)
	}
	if v, ok := d.GetOk("secondary_zone"); ok {
		pool.SecondaryZone = v.(string)
	}
	if v, ok := d.GetOk("billing_label"); ok {
		values := v.(*schema.Set)
		if values.Len() > 0 {
			labels := make([]billingLabel, 0, values.Len())
			for _, v := range expandBillingLabel(values) {
				labels = append(labels, v)
			}
			pool.BillingLabels = labels
		}
	}

	res, err := client.createStoragePool(&pool)
	if err != nil {
		log.Printf("Error creating storage pool: %#v", err)
		return err
	}
	d.SetId(res.PoolID)

	log.Printf("Created storage pool in region: %v", res.Region)

	return resourceGCPStoragePoolRead(d, meta)
}

func resourceGCPStoragePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	id := d.Id()
	pool := storagePool{}
	pool.Region = d.Get("region").(string)
	pool.PoolID = id
	var res storagePool
	res, err := client.getStoragePoolByID(&pool)
	if err != nil {
		return err
	}
	if res.PoolID != id {
		return fmt.Errorf("Expected storage pool with id: %v, Response contained storage pool with id: %v",
			d.Id(), res.PoolID)
	}
	if err := d.Set("size", res.SizeInBytes/GiBToBytes); err != nil {
		return fmt.Errorf("Error reading storage pool size: %s", err)
	}

	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("Error reading storage pool region: %s", err)
	}

	if err := d.Set("name", res.Name); err != nil {
		return fmt.Errorf("Error reading storage pool name: %s", err)
	}

	if _, ok := d.GetOk("billing_label"); ok {
		labels := flattenBillingLabel(res.BillingLabels)
		if err := d.Set("billing_label", labels); err != nil {
			return fmt.Errorf("Error reading storage pool billing_label: %s", err)
		}
	}

	return nil
}

func resourceGCPStoragePoolDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting storage pool: %#v", d.Get("name"))
	client := meta.(*Client)
	pool := storagePool{}
	pool.Region = d.Get("region").(string)
	pool.PoolID = d.Id()
	deleteErr := client.deleteStoragePool(&pool)
	if deleteErr != nil {
		return deleteErr
	}
	d.SetId("")

	return nil
}

func resourceGCPStoragePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	pool := storagePool{}
	// all of the following are required for API: update.
	pool.Region = d.Get("region").(string)
	pool.Name = d.Get("name").(string)
	pool.PoolID = d.Id()
	pool.ServiceLevel = d.Get("service_level").(string)

	if d.HasChange("size") {
		pool.SizeInBytes = d.Get("size").(int) * GiBToBytes
	}

	if d.HasChange("billing_label") {
		labels := d.Get("billing_label").(*schema.Set)
		pool.BillingLabels = expandBillingLabel(labels)
	}

	err := client.updateStoragePool(&pool)
	if err != nil {
		return err
	}
	return resourceGCPStoragePoolRead(d, meta)
}
