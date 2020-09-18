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

* `username` - (Required) The username of the Active Directory domain administrator.
* `password` - (Required) The password of the Active Directory domain administrator.
* `domain` - (Required) The name of the Active Directory domain.
* `region` - (Required) The region to which the Active Directory credentials are associated.
* `dns_server` - (Required) Comma separated list of DNS server IP addresses for the Active Directory domain.
* `net_bios` - (Required)  The netBIOS name prefix of the server.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the active directory.

## Unique id versus name

With NetApp_GCP, every resource has a unique id, but names are not necessarily unique.
