---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_volume"
sidebar_current: "docs-netapp-gcp-resource-volume"
description: |-
  Provides an NetApp_GCP volume resource. This can be used to create a new (empty) volume on the GCP-CVS.
---

# netapp_gcp\_volume

Provides an NetApp_GCP volume resource. This can be used to create a new (empty) volume on the GCP-CVS.

## Example Usages

**Read NetApp_GCP volume:**

```
data "netapp-gcp_volume" "data-volume" {
  name = "deleteme_asapGO_jusitin"
  region = "us-west2"
```

**Create NetApp_GCP volume:**

```
resource "netapp-gcp_volume" "gcp-volume" {
  provider = netapp-gcp
  name = "deleteme_asapGO_jusitin"
  region = "us-west2"
  protocol_types = ["NFSv3"]
  network = "cvs-terraform-vpc"
  size = 1024
  service_level = "premium"
  volume_path = "deleteme-asapGO"
  snapshot_policy {
    enabled = true
    daily_schedule {
      hour = 10
      minute = 1
    }
  }
  export_policy {
    rule {
      allowed_clients = "0.0.0.0/0"
      access= "ReadWrite"
      nfsv3 {
        checked =  true
      }
      nfsv4 {
        checked = false
      }
    }
    rule {
      allowed_clients= "10.0.0.0"
      access= "ReadWrite"
      nfsv3 {
        checked =  true
      }
      nfsv4 {
        checked = false
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `export_policy` - (Optional) The set of Export Policy attributes for volume.
* `name` - (Required) The name of the NetApp_GCP volume.
* `network` - (Required) The network VPC of the volume.
* `protocol_types` - (Required) The protocol_type of the volume. For NFS use 'NFSv3' or 'NFSv4' and for SMB use 'CIFS' or 'SMB'.
* `region` - (Required) The region where the NetApp_GCP volume to be created.
* `service_level` - (Optional) The performance of the service level of volume. Must be one of "standard", "premium", "extreme", default is "premium".
* `shared_vpc_project_number` - (Optional) The host project number when deploying in a shared VPC service project.
* `size` - (Required) The size of volume is between 1024 GiB to 102400 GiB inclusive.
* `snapshot_policy` - (Optional) The set of Snapshot Policy attributes for volume.
* `volume_path` - (Optional) The name of the volume path for volume.
* `type_dp` - (Optional) The type of the volume to be DP.
* `delete_on_creation_error` - (Optional) Delete volume if volume is in error state after creation. Default is false.
* `zone` - (Optional) The desired zone for the resource. If storage_class is set to 'software', zone is required.
* `storage_class` - (Optional) Storage Class to be provisioned. Allows the user to choose between hardware based or software based.
* `regional_ha` - (Optional) Flag indicating if the volume is regional, applicable only for software volumes.

The `snapshot_policy` block supports:
* `enabled` - (Optional) If enabled, make snapshots automatically according to the schedules. Default is false.
* `daily_schedule` - (Optional) If enabled, make a snapshot every day. Defaults to midnight.
* `hourly_schedule` - (Optional) If enabled, make a snapshot every hour e.g. at 04:00, 05:00, 06:00.
* `monthly_schedule` - (Optional) If enabled, make a snapshot every month at a specific day or days, defaults to the first day of the month at midnight
* `weekly_schedule` - (Optional) If enabled, make a snapshot every week at a specific day or days, defaults to Sunday at midnight.

The `daily_scheule` block supports:
* `hour` - (Optional) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional) The maximum number of Snapshots to keep for the daily schedule.

The `hourly_schedule` block supports:
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).

The `monthly_schedule` block supports:
* `days_of_month` - (Optional) Set the day or days of the month to make a snapshot (1-31). Accepts a comma delimited string of the day of the month e.g. '1,15,31'. Defaults to '1'.
* `hour` - (Optional) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional) The maximum number of Snapshots to keep for the daily schedule.

The `weekly_schedule` block supports:
* `day` - Set the day or days of the week to make a snapshot. Accepts a comma delimited string of week day names in english. Defaults to 'Sunday'.
* `hour` - (Optional) Set the hour to start the snapshot (0-23), defaults to midnight (0).
* `minute` - (Optional) Set the minute of the hour to start the snapshot (0-59), defaults to the top of the hour (0).
* `snapshots_to_keep` - (Optional) The maximum number of Snapshots to keep for the daily schedule.

The `export_policy` block supports:
* `rule` - (Optional) Export Policy rule.

The `rule` block supports:
* `access` - (Optional) Defines the access type for clients matching the 'allowedClients' specification.
* `allowed_clients` - (Optional) Defines the client ingress specification (allowed clients) as a comma seperated string with IPv4 CIDRs, IPv4 host addresses and host names.
* `has_root_access` - (Optional) If enabled (true or on) the rule defines that no_root_squash is set, else if it is disable (false or off) root_squash is set and user ID mapped to anonymous user.
* `nfsv3` - (Optional) If enabled (true) the rule allows NFSv3 protocol for clients matching the 'allowedClients' specification.
* `nfsv4` - (Optional) If enabled (true) the rule allows NFSv4 protocol for clients matching the 'allowedClients' specification.

The `nfsv3` block supports:
* `checked` - (Optional) Enable NFSv3 protocol.

The `nfsv4` block supports:
* `checked` - (Optional) Enable NFSv4 protocol.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the volume.

## Unique id versus name

With NetApp_GCP, every resource has a unique id, but names are not necessarily unique.
