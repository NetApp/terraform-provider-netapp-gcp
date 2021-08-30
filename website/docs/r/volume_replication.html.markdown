---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_volume_replication"
sidebar_current: "docs-netapp-gcp-resource-volume-replication"
description: |-
  Provides an NetApp_GCP volume resource. This can be used to create a new volume replication relationship on the GCP-CVS.
---

# netapp_gcp\_volume\_replication

Provides an NetApp_GCP volume replication resource. This can be used to create a new volume replication relationship on the GCP-CVS.

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

* `bandwidth` - (Optional) Maximum bandwidth that the replication will allow.
* `destination_volume_uuid` - (Optional) UUID v4 of the destination volume of a volume replication relationship.
* `endpoint_type` - (Optional) Indicates whether the local volume is the source or destination for the Volume Replication.
* `name` - (Required) The name of the NetApp_GCP volume replication.
* `policy` - (Optional) Replication policy.
* `region` - (Required) The region where the NetApp_GCP volume replication relationship to be created.
* `remote_region` - (Optional) The remote region for the other end of the volume replication.
* `schedule` - (Optional) Replication_policy.
* `source_volume_uuid` - (Optional) UUID v4 of the source volume of a volume replication relationship.


