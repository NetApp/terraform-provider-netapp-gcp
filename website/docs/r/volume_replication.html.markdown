---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_volume_replication"
sidebar_current: "docs-netapp-gcp-resource-volume-replication"
description: |-
  Provides an NetApp_GCP volume resource. This can be used to create a new volume replication relationship on the GCP-CVS.
---

# netapp_gcp\_volume\_replication

Provides an NetApp_GCP volume replication resource. This can be used to create a new volume replication relationship on the GCP-CVS.

Destinations pull from source. `source_volume_uuid` in `remote_region` is replicated to `destination_volume_id` in `region`. 

## Example Usages

**Create NetApp_GCP volume:**

```
resource "netapp-gcp_volume_replication" "gcp-volume-replication" {
  provider = netapp-gcp
  name = "myReplica"
  region = "us-east4"
  remote_region = "us-east2"
  destination_volume_id = "12345678-abcd-abcd-abcd-123456789012"
  source_volume_id = "87654321-abcd-abcd-abcd-123456789012"
  endpoint_type = "dst"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp_GCP volume replication.
* `source_volume_uuid` - (Required) UUID v4 of the destination volume of a volume replication relationship.
* `remote_region` - (Required) The region of the source volume.
* `destination_volume_id` - (Required) UUID v4 of the source volume of a volume replication relationship.
* `region` - (Required) The region of the destination volume.
* `endpoint_type` - (Required) Always set "dst".
* `schedule` - (Required) Replication_policy ("10minutely", "hourly", "daily")
* `policy` - (Optional) Replication policy.
