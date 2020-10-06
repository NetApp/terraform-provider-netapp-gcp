---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_snapshot"
sidebar_current: "docs-netapp-gcp-resource-snapshot"
description: |-
  Provides a NetApp_GCP snapshot resource. This can be used to create a new snapshot on the CVS for GCP.
---

# netapp_gcp\_snapshot

Provides a NetApp_GCP snapshot resource. This can be used to create a new snapshot on the CVS for GCP.

## Example Usages

**Create NetApp_GCP snapshot:**

```
resource "netapp-gcp_snapshot" "gcp-snapshot" {
  name = "main-snapshot"
  region = "us-west2"
  volume_name =  "main-volume"
  creation_token = "unique-token-number"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp_GCP snapshot to be created.
* `region` - (Required) The region where the NetApp_GCP volume exists.
* `volume_name` - (Optional) The name of the volume to create a snapshot from.
* `creation_token` - (Optional) The creation token of volume of the NetApp_GCP.

 At least one of volume_name or creation_token is required to create snapshot.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the snapshot.

## Unique id versus name

With NetApp_GCP, every resource has a unique id, but names are not necessarily unique. Make sure that volume names are unique within a region for a given subscription when Creation Token parameter is not used.
