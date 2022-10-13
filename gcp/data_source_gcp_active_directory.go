package gcp

import (
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
			"net_bios": {
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
			"aes_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backup_operators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"security_operators": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"allow_local_nfs_users_with_ldap": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"kdc_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_signing": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ad_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed_ad": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceGCPActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	return resourceGCPActiveDirectoryRead(d, meta)
}
