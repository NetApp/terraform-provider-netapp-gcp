package gcp

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp/cvs/restapi"
)

func resourceGCPActiveDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPActiveDirectoryCreate,
		Read:   resourceGCPActiveDirectoryRead,
		Delete: resourceGCPActiveDirectoryDelete,
		Exists: resourceGCPActiveDirectoryExists,
		Update: resourceGCPActiveDirectoryUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			//these available fields are required for create and update.
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"net_bios": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organizational_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGCPActiveDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating active directory: %#v", d)
	client := meta.(*Client)
	// check whether the AD already exists on GCP, if it exist, error out.
	listActiveDirectory := listActiveDirectoryRequest{}
	listActiveDirectory.Region = d.Get("region").(string)
	existedAd, err := client.listActiveDirectoryForRegion(listActiveDirectory)
	if err != nil {
		log.Print("Error checking current active directory before creating new active directory.")
		return err
	}
	if existedAd.UUID != "" {
		return fmt.Errorf("Active Directory in region: \"%v\" already exists", existedAd.Region)
	}

	activeDirectory := operateActiveDirectoryRequest{}
	activeDirectory.Username = d.Get("username").(string)
	activeDirectory.Password = d.Get("password").(string)
	activeDirectory.Domain = d.Get("domain").(string)
	activeDirectory.DNS = d.Get("dns_server").(string)
	activeDirectory.NetBIOS = d.Get("net_bios").(string)
	if v, ok := d.GetOk("organizational_unit"); ok {
		activeDirectory.OrganizationalUnit = v.(string)
	}
	if v, ok := d.GetOk("site"); ok {
		activeDirectory.Site = v.(string)
	}
	activeDirectory.Region = d.Get("region").(string)

	res, err := client.createActiveDirectory(&activeDirectory)
	if err != nil {
		log.Print("Error creating active directory")
		return err
	}
	d.SetId(res.UUID)

	log.Printf("Created active directory in region: %v", activeDirectory.Region)

	return resourceGCPActiveDirectoryRead(d, meta)
}

func resourceGCPActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	id := d.Id()
	activeDirectory := listActiveDirectoryRequest{}
	activeDirectory.Region = d.Get("region").(string)
	var res listActiveDirectoryResult
	res, err := client.listActiveDirectoryForRegion(activeDirectory)
	if err != nil {
		return err
	}
	if res.UUID != id {
		return fmt.Errorf("Expected active directory with id: %v, Response contained active directory with id: %v",
			d.Get("uuid").(string), res.UUID)
	}
	d.Set("uuid", res.UUID)

	if err := d.Set("domain", res.Domain); err != nil {
		return fmt.Errorf("Error reading active directory domain: %s", err)
	}

	if err := d.Set("net_bios", res.NetBIOS); err != nil {
		return fmt.Errorf("Error reading active directory net_bios: %s", err)
	}

	if err := d.Set("organizational_unit", res.OrganizationalUnit); err != nil {
		return fmt.Errorf("Error reading active directory organizational_unit: %s", err)
	}

	if err := d.Set("site", res.Site); err != nil {
		return fmt.Errorf("Error reading active directory site: %s", err)
	}

	if err := d.Set("username", res.Username); err != nil {
		return fmt.Errorf("Error reading active directory username: %s", err)
	}

	if err := d.Set("dns_server", res.DNS); err != nil {
		return fmt.Errorf("Error reading active directory dns_server: %s", err)
	}

	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("Error reading active directory region: %s", err)
	}

	return nil
}

func resourceGCPActiveDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting active directory: %#v", d)
	client := meta.(*Client)
	activeDirectory := deleteActiveDirectoryRequest{}
	activeDirectory.Region = d.Get("region").(string)
	activeDirectory.UUID = d.Get("uuid").(string)
	deleteErr := client.deleteActiveDirectory(activeDirectory)
	if deleteErr != nil {
		return deleteErr
	}
	d.SetId("")

	return nil
}

func resourceGCPActiveDirectoryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of active directory: %#v", d)
	client := meta.(*Client)
	activeDirectory := listActiveDirectoryRequest{}
	activeDirectory.UUID = d.Get("uuid").(string)
	activeDirectory.Region = d.Get("region").(string)
	var res listActiveDirectoryResult
	res, err := client.listActiveDirectoryForRegion(activeDirectory)
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
	if res.UUID != activeDirectory.UUID {
		d.SetId("")
		return false, nil
	}

	return true, err
}

func resourceGCPActiveDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Checking existence of active directory: %#v", d)
	client := meta.(*Client)
	activeDirectory := operateActiveDirectoryRequest{}
	// all of the following are required for API: update.
	activeDirectory.Username = d.Get("username").(string)
	activeDirectory.Password = d.Get("password").(string)
	activeDirectory.Domain = d.Get("domain").(string)
	activeDirectory.DNS = d.Get("dns_server").(string)
	activeDirectory.NetBIOS = d.Get("net_bios").(string)
	activeDirectory.OrganizationalUnit = d.Get("organizational_unit").(string)
	activeDirectory.Site = d.Get("site").(string)
	activeDirectory.Region = d.Get("region").(string)
	activeDirectory.UUID = d.Get("uuid").(string)
	err := client.updateActiveDirectory(activeDirectory)
	if err != nil {
		return err
	}
	return resourceGCPActiveDirectoryRead(d, meta)
}
