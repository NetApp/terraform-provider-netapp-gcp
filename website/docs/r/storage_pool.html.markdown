---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_storage_pool"
sidebar_current: "docs-netapp-gcp-resource-storage-pool"
description: |-
  Provides a NetApp_GCP storage pool resource. This can be used to create a new storage pool on the GCP-CVS.
---

# netapp_gcp\_storage_pool

Provides a NetApp_GCP storage pool resource. This can be used to manage storage pools for volumes of service-type CVS.

## Example Usages

**Create NetApp_GCP zonal storage pool:**

```
resource "netapp-gcp_storage_pool" "test-storage-pool" {
  name = "example_pool"
  region = "us-east1"
  zone = "us-east1-b"
  network = "example-vpc"
  global_ad_access = false
  size = 1024
  service_level = "StandardSW"
  storage_class = "software"
  billing_label {
      key = "exampleKey"
      value = "exampleValue"
    }
}
```

**Create NetApp_GCP regional storage pool:**

```
resource "netapp-gcp_storage_pool" "test-storage-pool" {
  name = "example_ha_pool"
  region = "us-east1"
  zone = "us-east1-b"
  secondary_zone = "us-east1-c"
  network = "example-vpc"
  global_ad_access = true
  size = 1024
  service_level = "ZoneRedundantStandardSW"
  storage_class = "software"
  billing_label {
      key = "exampleKey"
      value = "exampleValue"
    }
}
```

**Create NetApp_GCP volume in these pools:**

Volumes in pools require some pool specific arguments. The example below shows how to supply these arguments.
```
resource "netapp-gcp_volume" "gcp-volume-nfs" {
  # This block mirrors pool settings
  region = netapp-gcp_storage_pool.regional-storage-pool.region
  zone = netapp-gcp_storage_pool.regional-storage-pool.zone
  network = netapp-gcp_storage_pool.regional-storage-pool.network
  pool_id = netapp-gcp_storage_pool.regional-storage-pool.id
  regional_ha = regional_ha = netapp-gcp_storage_pool.regional-storage-pool.service_level == "ZoneRedundantStandardSW" ? true : false
  storage_class = "software" 
  service_level = "standard"  # leave as is

  # Volume specific settings
  name = "example_volume"
  protocol_types = ["NFSv3"]
  size = 200

  # Snapshot policy definition
  snapshot_policy {
    enabled = true
  }
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the storage pool.
* `region` - (Required) The region where the storage pool to be created.
* `zone` - (Required) Location of the pool.
* `size` - (Required, modifiable) Storage pool size.
* `network` - (Required) Network name.
* `global_ad_access` - (Optional, modifiable) Enables global access to Active Directory controllers outside of the pools region.
* `service_level` - (Required) StandardSW or ZoneRedundantStandardSW.
* `storage_class` - (Required) Software.
* `billing_label` - (Optional, modifiable) Key-value pair for billing labels.
* `shared_vpc_project_number` - (Optional) The host project number when deploying in a shared VPC service project.
* `regional_ha` - (Optional, deprecated) Flag indicating if the pool is regional, applicable only for software type. Is replaced by service_level.
* `secondary_zone` - (Optional) Secondary zone if service level is ZoneRedundantStandardSW.

The `billing_label` block supports:
* `key` - (Required) Must be a minimum length of 1 character and a maximum length of 63 characters, and cannot be empty. Can contain only lowercase letters, numeric characters, underscores, and dashes. All characters must use UTF-8 encoding, and international characters are allowed. Must start with a lowercase letter or international character.
* `value` - (Required) Can be empty, and have a maximum length of 63 characters. Can contain only lowercase letters, numeric characters, underscores, and dashes. All characters must use UTF-8 encoding, and international characters are allowed.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the storage pool.
* `managed_pool` - A pool which was automatically created when using creating pre-StoragePool volumes. See [Managed Pools](https://cloud.google.com/architecture/partners/netapp-cloud-volumes/storage-pools?hl=en_US#managed_pools)
