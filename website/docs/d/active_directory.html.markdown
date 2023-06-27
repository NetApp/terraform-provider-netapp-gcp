---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_active_directory"
sidebar_current: "docs-netapp-gcp-data-source-active-directory"
description: |-
  Provides an NetApp_GCP active directory resource. This can be used to get a existing active directory on the GCP-CVS.
---

# netapp_gcp\_active\_directory

Provides an NetApp_GCP active directory resource. This can be used to get a existing active directory on the GCP-CVS.

## Example Usages

**Read NetApp_GCP active directory:**

Query data from an existing Active Directory connection for a given region.

```
data "netapp-gcp_active_directory" "ad-us-west2" {
    region = "us-west2"
}
```

## Argument Reference

The following arguments are supported:

AD connection specific settings:
* `region` - (Required) The region to which the Active Directory credentials are associated.

  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the active directory.
* `domain` - Fully qualified name of Active Directory domain.
* `net_bios` - NetBIOS prefix name of the server that will be created. A random 5-digit suffix is appended automatically (e.g. -A579).
* `organizational_unit` - Name of the Organizational Unit (OU) within Windows Active Directory the computer account belongs to in order from leaf OU to root OU. Defaults to CN=Computers if left empty.
* `site` -  Specify an Active Directory site to manage domain controller selection. Use when Active Directory domain controllers in multiple regions are configured. Defaults to “Default-First-Site-Name” if left empty.
* `dns_server` - Comma separated list of DNS server IP addresses used for DNS-based domain controller discovery.
* `aes_encryption` - Enables AES-128 and AES-256 encryption for Kerberos-based communication with Active Directory. Default is false.
* `ldap_signing` - Enables LDAP siging. Default is false.
* `connection_type` - Specify "software" for service type CVS or "hardware" for service type CVS-Performance.
* `managed_ad` - Flags this configuration as Google ManagedAD configuration. Please see https://cloud.google.com/architecture/partners/netapp-cloud-volumes/managing-active-directory-connections?hl=en_US#connect_to_managed_microsoft_ad


User credentials for Domain join:
* `username` - Username of an account permitted to create computer objects in your Active Directory.

SMB specific settings:
* `backup_operators` - Users to be added to the built-in Backup Operators active directory group. The usernames must be unique, and entries cannot include @ or \\. The entire list will be validated and rejected as whole if one or more entries are invalid.
* `security_operators` - Domain users to be given the SeSecurityPrivilege.

NFS specific settings:
* `ad_server` - Hostname of an Active Directory domain controller which is used as Kerberos Key Distribution Center (KDC). This optional parameter is used only for Kerberized NFS.
* `kdc_ip` - IP address of the Active Directory domain controller which is used as Kerberos Key Distribution Center (KDC). This optional parameter is used only for Kerberized NFS.
* `allow_local_nfs_users_with_ldap` - If enabled, allow_local_nfs_users_with_ldap will allow access to local users as well as LDAP users. If access is needed for only LDAP users, it has to be disabled. Default is false.
