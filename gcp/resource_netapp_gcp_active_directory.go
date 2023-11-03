package gcp

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"aes_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				Default:  false,
			},
			"kdc_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_signing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connection_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"software", "hardware"}, true),
			},
			"ad_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed_ad": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	activeDirectory.Label = d.Get("connection_type").(string)
	if v, ok := d.GetOk("organizational_unit"); ok {
		activeDirectory.OrganizationalUnit = v.(string)
	}
	if v, ok := d.GetOk("site"); ok {
		activeDirectory.Site = v.(string)
	}
	activeDirectory.Region = d.Get("region").(string)

	activeDirectory.AesEncryption = d.Get("aes_encryption").(bool)
	activeDirectory.LdapSigning = d.Get("ldap_signing").(bool)
	activeDirectory.AllowLocalNFSUsersWithLdap = d.Get("allow_local_nfs_users_with_ldap").(bool)

	if v, ok := d.GetOk("backup_operators"); ok {
		backupOperators := make([]string, 0)
		for _, y := range v.(*schema.Set).List() {
			backupOperators = append(backupOperators, y.(string))
		}
		activeDirectory.BackupOperators = backupOperators
	}

	if v, ok := d.GetOk("security_operators"); ok {
		securityOperators := make([]string, 0)
		for _, y := range v.(*schema.Set).List() {
			securityOperators = append(securityOperators, y.(string))
		}
		activeDirectory.SecurityOperators = securityOperators
	}

	if v, ok := d.GetOk("kdc_ip"); ok {
		activeDirectory.KdcIP = v.(string)
	}

	if v, ok := d.GetOk("ad_server"); ok {
		activeDirectory.AdName = v.(string)
	}

	activeDirectory.ManagedAD = d.Get("managed_ad").(bool)

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
	activeDirectory := listActiveDirectoryRequest{}
	activeDirectory.Region = d.Get("region").(string)
	activeDirectory.UUID = d.Id()
	var res listActiveDirectoryResult
	res, err := client.listActiveDirectoryForRegion(activeDirectory)
	if err != nil {
		return err
	}
	// Disabling, since it would fail for call from dataSourceGCPVolumeRead
	// Unclear if this sanity check is required
	// if res.UUID != d.id {
	// 	return fmt.Errorf("Expected active directory with id: %v, Response contained active directory with id: %v",
	// 		d.Get("uuid").(string), res.UUID)
	// }
	d.SetId(res.UUID)
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

	if err := d.Set("aes_encryption", res.AesEncryption); err != nil {
		return fmt.Errorf("Error reading active directory aes_encryption: %s", err)
	}

	if err := d.Set("ldap_signing", res.LdapSigning); err != nil {
		return fmt.Errorf("Error reading active directory ldap_signing: %s", err)
	}

	if err := d.Set("allow_local_nfs_users_with_ldap", res.AllowLocalNFSUsersWithLdap); err != nil {
		return fmt.Errorf("Error reading active directory allow_local_nfs_users_with_ldap: %s", err)
	}

	if err := d.Set("security_operators", res.SecurityOperators); err != nil {
		return fmt.Errorf("Error reading active directory security_operators: %s", err)
	}

	if err := d.Set("backup_operators", res.BackupOperators); err != nil {
		return fmt.Errorf("Error reading active directory backup_operators: %s", err)
	}

	if err := d.Set("kdc_ip", res.KdcIP); err != nil {
		return fmt.Errorf("Error reading active directory kdc_ip: %s", err)
	}

	if err := d.Set("connection_type", res.Label); err != nil {
		return fmt.Errorf("Error reading active directory connection_type: %s", err)
	}

	if err := d.Set("ad_server", res.AdName); err != nil {
		return fmt.Errorf("Error reading active directory ad_server: %s", err)
	}

	if err := d.Set("managed_ad", res.ManagedAD); err != nil {
		return fmt.Errorf("Error reading active directory managed_ad: %s", err)
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
	activeDirectory.UUID = d.Id()
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
	activeDirectory.Label = d.Get("connection_type").(string)

	activeDirectory.AesEncryption = d.Get("aes_encryption").(bool)

	if v, ok := d.GetOk("backup_operators"); ok {
		backupOperators := make([]string, 0)
		for _, y := range v.(*schema.Set).List() {
			backupOperators = append(backupOperators, y.(string))
		}
		activeDirectory.BackupOperators = backupOperators
	}

	if v, ok := d.GetOk("security_operators"); ok {
		securityOperators := make([]string, 0)
		for _, y := range v.(*schema.Set).List() {
			securityOperators = append(securityOperators, y.(string))
		}
		activeDirectory.SecurityOperators = securityOperators
	}

	activeDirectory.AllowLocalNFSUsersWithLdap = d.Get("allow_local_nfs_users_with_ldap").(bool)

	if d.HasChange("kdc_ip") {
		activeDirectory.KdcIP = d.Get("kdc_ip").(string)
	}

	activeDirectory.LdapSigning = d.Get("ldap_signing").(bool)

	if v, ok := d.GetOk("ad_server"); ok {
		activeDirectory.AdName = v.(string)
	}

	activeDirectory.ManagedAD = d.Get("managed_ad").(bool)

	err := client.updateActiveDirectory(activeDirectory)
	if err != nil {
		return err
	}
	return resourceGCPActiveDirectoryRead(d, meta)
}
