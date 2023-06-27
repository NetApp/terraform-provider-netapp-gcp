---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_volume"
sidebar_current: "docs-netapp-gcp-data-source-volume"
description: |-
  Provides an NetApp_GCP volume data source. This can be used to get a existing volume on the GCP-CVS.
---

# netapp_gcp\_volume

Provides an NetApp_GCP volume resource. This can be used to manage volumes on GCP-CVS.

## Example Usages

**Read existing NetApp_GCP volume:**

```
data "netapp-gcp_volume" "data-volume" {
  name = "deleteme_asap"
  region = "us-west2"
```


## Argument Reference

The following arguments are supported:

Generic volume settings
* `region` - (Required) The region where the NetApp_GCP volume to be created.
* `name` - (Required) The name of the NetApp_GCP volume.


## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `id` - The unique identifier for the volume.

