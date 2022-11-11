package gcp

import (
	"fmt"
	"log"
	"regexp"
	"strings"

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
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Please use service_level = StandardSW or ZoneRedundantStandardSW instead",
			},
			"global_ilb": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"managed_pool": {
				Type:     schema.TypeBool,
				Computed: true,
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
	if v, ok := d.GetOk("global_ilb"); ok {
		pool.GlobalILB = v.(bool)
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
			labels = append(labels, expandBillingLabel(values)...)
			pool.BillingLabels = labels
		}
	}

	if v, ok := d.GetOk("shared_vpc_project_number"); ok {
		pool.SharedVpcProjectNumber = v.(string)
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
	if strings.ToLower(res.State) == "deleted" {
		d.SetId("")
		return nil
	}
	if res.PoolID != id {
		return fmt.Errorf("expected storage pool with id: %v, Response contained storage pool with id: %v",
			d.Id(), res.PoolID)
	}
	if err := d.Set("size", res.SizeInBytes/GiBToBytes); err != nil {
		return fmt.Errorf("error reading storage pool size: %s", err)
	}

	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("error reading storage pool region: %s", err)
	}

	if err := d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error reading storage pool name: %s", err)
	}

	if _, ok := d.GetOk("billing_label"); ok {
		labels := flattenBillingLabel(res.BillingLabels)
		if err := d.Set("billing_label", labels); err != nil {
			return fmt.Errorf("error reading storage pool billing_label: %s", err)
		}
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
				return fmt.Errorf("error reading shared_vpc_project_number: %s", err)
			}
		}
	} else {
		return fmt.Errorf("network path %s is invalid", res.Network)
	}
	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("error reading volume network: %s", err)
	}

	if err := d.Set("global_ilb", res.GlobalILB); err != nil {
		return fmt.Errorf("error reading storage pool global_ilb flag: %s", err)
	}

	if err := d.Set("managed_pool", res.ManagedPool); err != nil {
		return fmt.Errorf("error reading storage pool managed_pool flag: %s", err)
	}

	if err := d.Set("zone", res.Zone); err != nil {
		return fmt.Errorf("error reading storage pool zone: %s", err)
	}

	if err := d.Set("secondary_zone", res.SecondaryZone); err != nil {
		return fmt.Errorf("error reading storage pool secondary_zone: %s", err)
	}

	if err := d.Set("service_level", res.ServiceLevel); err != nil {
		return fmt.Errorf("error reading storage pool service_level: %s", err)
	}
	// RegionalHA is old parameter used in legacy volumes (managed_pool)
	// It should not be used anymore and is replaced by
	// service_level = StandardSW or ZoneRedundantStandardSW
	//
	// Calculate internal state for virtual parameter RegionalHA
	res.RegionalHA = false
	if res.ServiceLevel == "ZoneRedundantStandardSW" {
		res.RegionalHA = true
	}
	// if err := d.Set("regional_ha", res.RegionalHA); err != nil {
	// 	return fmt.Errorf("error setting storage pool regional_ha: %s", err)
	// }

	if err := d.Set("storage_class", res.StorageClass); err != nil {
		return fmt.Errorf("error reading storage pool storage_class: %s", err)
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

	if d.HasChange("global_ilb") {
		pool.GlobalILB = d.Get("global_ilb").(bool)
	}

	if d.HasChange("zone") {
		pool.Zone = d.Get("zone").(string)
	}

	err := client.updateStoragePool(&pool)
	if err != nil {
		return err
	}
	return resourceGCPStoragePoolRead(d, meta)
}
