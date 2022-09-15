---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_active_directory"
sidebar_current: "docs-netapp-gcp-resource-active-directory"
description: |-
  Provides an NetApp_GCP active directory resource. This can be used to create a new active directory on the GCP-CVS.
---

# netapp_gcp\_active\_directory

Provides an NetApp_GCP active directory resource. This can be used to create a new active directory on the GCP-CVS.

## Example Usages

**Read NetApp_GCP active directory:**

Query data from an existing Active Directory connection for a given region.

```
data "netapp-gcp_active_directory" "ad-us-west2" {
    region = "us-west2"
}
```

**Create NetApp_GCP active directory:**
Creates a new Active Directory connection for a given region. Only one connection per region and connection_type allowed.

```
resource "netapp-gcp_active_directory" "gcp-active-directory" {
  provider = netapp-gcp
  region = "us-west2"
	username = "test_user"
	password = "netapp"
  domain = "example.com"
  dns_server = "10.0.0.0"
  net_bios = "cvsserver"
  connection_type = "hardware"
}
```

## Argument Reference

The following arguments are supported:

AD connection specific settings:
* `region` - (Required) The region to which the Active Directory credentials are associated.
* `connection_type` - (Required) Specify "software" for service type CVS or "hardware" for service type CVS-Performance.
* `domain` - (Required) Fully qualified name of Active Directory domain.
* `dns_server` - (Required) Comma separated list of DNS server IP addresses used for DNS-based domain controller discovery..
* `site` - (Optional) Specify an Active Directory site to manage domain controller selection. Use when Active Directory domain controllers in multiple regions are configured. Defaults to “Default-First-Site-Name” if left empty.
* `organizational_unit` - (Optional) Name of the Organizational Unit (OU) within Windows Active Directory the computer account belongs to in order from leaf OU to root OU. Defaults to CN=Computers if left empty.
* `net_bios` - (Required) NetBIOS prefix name of the server that will be created. A random 5-digit suffix is appended automatically (e.g. -A579).
* `aes_encryption` - (Optional) Enables AES-128 and AES-256 encryption for Kerberos-based communication with Active Directory. Default is false.
* `ldap_signing` - (Optional) Enables LDAP siging. Default is false.
* `managed_ad` - (Optional) Flags this configuration as Google ManagedAD configuration. Please see https://cloud.google.com/architecture/partners/netapp-cloud-volumes/managing-active-directory-connections?hl=en_US#connect_to_managed_microsoft_ad

User credentials for Domain join:
* `username` - (Required) Username of an account permitted to create computer objects in your Active Directory.
* `password` - (Required) Password for the account permitted to create computer objects in your Active Directory.

SMB specific settings:
* `backup_operators` - (Optional) Users to be added to the built-in Backup Operators active directory group. The usernames must be unique, and entries cannot include @ or \\. The entire list will be validated and rejected as whole if one or more entries are invalid.
* `security_operators` - (Optional) Domain users to be given the SeSecurityPrivilege.

NFS specific settings:
* `ad_server` - (Optional) Hostname of an Active Directory domain controller which is used as Kerberos Key Distribution Center (KDC). This optional parameter is used only for Kerberized NFS.
* `kdc_ip` - (Optional) IP address of the Active Directory domain controller which is used as Kerberos Key Distribution Center (KDC). This optional parameter is used only for Kerberized NFS.
* `allow_local_nfs_users_with_ldap` - (Optional) If enabled, allow_local_nfs_users_with_ldap will allow access to local users as well as LDAP users. If access is needed for only LDAP users, it has to be disabled. Default is false.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the active directory.

## Unique id versus name

With NetApp_GCP, every resource has a unique id, but names are not necessarily unique.