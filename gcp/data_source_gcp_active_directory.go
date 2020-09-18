package gcp

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGCPActiveDirectory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGCPActiveDirectoryRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"netbios": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceGCPActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	active_directory := listActiveDirectoryRequest{}
	active_directory.Region = d.Get("region").(string)
	var res listActiveDirectoryResult
	res, err := client.listActiveDirectoryForRegion(active_directory)
	if err != nil {
		return err
	}
	d.SetId(res.UUID)

	if err := d.Set("uuid", res.UUID); err != nil {
		return fmt.Errorf("Error reading active directory UUID: %s", err)
	}
	if err := d.Set("domain", res.Domain); err != nil {
		return fmt.Errorf("Error reading active directory domain: %s", err)
	}

	if err := d.Set("netbios", res.NetBIOS); err != nil {
		return fmt.Errorf("Error reading active directory netbios: %s", err)
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
