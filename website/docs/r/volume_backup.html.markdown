---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_volume_backup"
sidebar_current: "docs-netapp-gcp-resource-volume-backup"
description: |-
  Provides a NetApp_GCP volume_backup resource. This can be used to create a new volume_backup for required volume on the CVS for GCP.
---

# netapp_gcp\_volume_backup

Provides a NetApp_GCP volume_backup resource. This can be used to create a new volume_backup for required volume on the CVS for GCP.

## Example Usages

**Create NetApp_GCP volume_backup:**

```
resource "netapp-gcp_volume_backup" "gcp-volume-backup" {
  name = "main-volume-backup"
  region = "us-west2"
  volume_name =  "main-volume"
  creation_token = "unique-token-number"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp_GCP volume_backup to be created.
* `region` - (Required) The region where the NetApp_GCP volume exists.
* `volume_name` - (Required) The name of the volume to create a volume_backup from.
* `creation_token` - (Required) The creation token of volume of the NetApp_GCP.

 At least one of volume_name or creation_token is required to create volume_backup.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the volume_backup.

## Unique id versus name

With NetApp_GCP, every resource has a unique id, but names are not necessarily unique. Make sure that volume names are unique within a region for a given subscription when Creation Token parameter is not used.
