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

```
data "netapp-gcp_active_directory" "ad-us-west2" {
    region = "us-west2"
}
```

**Create NetApp_GCP active directory:**

```
resource "netapp-gcp_active_directory" "gcp-active-directory" {
  provider = netapp-gcp
  region = "us-west2"
	username = "test_user"
	password = "netapp"
  domain = "example.com"
  dns_server = "10.0.0.0"
  net_bios = "cvsserver"
}
```

## Argument Reference

The following arguments are supported:

* `aes_encryption` - (Optional) If enabled, AES encryption will be enabled for SMB communication. Default is false.
* `allow_local_nfs_users_with_ldap` - (Optional) If enabled, allow_local_nfs_users_with_ldap will allow access to local users as well as LDAP users. If access is needed for only LDAP users, it has to be disabled. Default is false.
* `backup_operators` - (Optional) Users to be added to the Built-in Backup Operator active directory group. The usernames must be unique, and entries cannot include @ or \\. The entire list will be validated and rejected as whole if one or more entries are invalid.
* `dns_server` - (Required) Comma separated list of DNS server IP addresses for the Active Directory domain.
* `domain` - (Required) The name of the Active Directory domain.
* `kdc_ip` - (Optional)  kdc server IP address for the active directory machine. This optional parameter is used only while creating kerberos volume.
* `ldap_signing` - (Optional) Specifies whether or not the LDAP traffic needs to be signed. Default is false.
* `net_bios` - (Required) The netBIOS name prefix of the server.
* `password` - (Required) The password of the Active Directory domain administrator.
* `region` - (Required) The region to which the Active Directory credentials are associated.
* `security_operators` - (Optional) Domain users to be given the SeSecurityPrivilege.
* `username` - (Required) The username of the Active Directory domain administrator.

  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the active directory.

## Unique id versus name

With NetApp_GCP, every resource has a unique id, but names are not necessarily unique.
