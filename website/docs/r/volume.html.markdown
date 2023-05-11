---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_volume"
sidebar_current: "docs-netapp-gcp-resource-volume"
description: |-
  Provides an NetApp_GCP volume resource. This can be used to create a new (empty) volume on the GCP-CVS.
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

**Create NetApp_GCP volume:**

```
resource "netapp-gcp_volume" "gcp-volume" {
  provider = netapp-gcp
  name = "deleteme_asap"
  region = "us-west2"
  protocol_types = ["NFSv3"]
  network = "cvs-terraform-vpc"
  size = 1024
  service_level = "premium"
  volume_path = "deleteme-asap"
  snapshot_policy {
    enabled = true
    daily_schedule {
      snapshotsToKeep = 7
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

Generic volume settings
* `region` - (Required) The region where the NetApp_GCP volume to be created.

* `name` - (Required) The name of the NetApp_GCP volume.
* `volume_path` - (Optional) The name of the export path or share name to be used for the volume. Must be unique per region.
* `shared_vpc_project_number` - (Optional) The host project number when deploying in a shared VPC service project.
* `network` - (Required) Name of VPC network for the volume.
* `size` - (Required) The size of volume. 100-102400 GiB for CVS-Performance, 1-102400 GiB for CVS on Storage Pools
* `delete_on_creation_error` - (Optional) Automatically delete volume if volume is in error state after creation. Default is false.

Service-Type CVS-Performance specific settings:
* `storage_class` - "hardware" for CVS-Performance.
* `service_level` - (Optional) The performance of the service level of volume. Must be one of "standard", "premium", "extreme", default is "premium".
* `type_dp` - (Optional) True for Volume Replication destination volume, False for normal primary volume.
* `snapshot_id` - (Optional) The UUID of the snapshot to create volume from.

Service-Type CVS specific settings:
* `storage_class` - "software" for CVS.
* `service_level` - "standard" for CVS.
* `pool_id` - (Required) UUID of the PoolId under which volumes get created.
* `regional_ha` - (Required) Flag indicating if the volume is regional, applicable only for software volumes. Set true for zone redundant Storage Pools.
* `zone` - (Required) The desired zone for the resource. If storage_class is set to 'software', zone is required. 

Protocol settings:
* `protocol_types` - (Required) The protocol_type of the volume. For CVS use 'NFSv3' or 'SMB'. For CVS-Performance use 'NFSv3', 'NFSv4' or 'SMB', or a combinations of ['NFSv3', 'NFSv4'], ['SMB', 'NFSv3'] and ['SMB', 'NFSv4'].
* `security_style` - (Optional) Security style for dual-protocol volumes of protocol_type ['SMB', 'NFSv3'] or ['SMB', 'NFSv4']. Valid choices are 'unix' and 'ntfs'. Pure NFS volumes will always be unix, pure SMB volumes always be ntfs. Setting cannot be changed after volume creation.

SMB protocol specific settings:
* `smb_share_settings` - (Optional) List of SMB share properties. Must be zero or more of "encrypt_data", "browsable", "changenotify", "non_browsable", "oplocks", "showsnapshot", "show_previous_versions", "continuously_available", "access_based_enumeration".

NFS protocol specific settings:
* `export_policy` - (Optional) Specify NFS Export Policy.
* `unix_permissions` - (Optional) UNIX permissions for root directory of NFS volume. Accepts octal 4 digit format. First digit selects the set user ID(4), set group ID (2) and sticky (1) attributes. Second digit selects permission for the owner of the file: read (4), write (2) and execute (1). Third selects permissions for other users in the same group. the fourth for other users not in the group. "0755" - gives read/write/execute permissions to owner and read/execute to group and other users.

The `export_policy` block (for NFSv3 and NFSv4) supports:
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

Snapshot settings:
* `snapshot_policy` - (Optional) The set of Snapshot Policy attributes for volume.
* `snapshot_directory` - (Optional) If enabled, the volume will contain a read-only .snapshot directory which provides access to each of the volume's snapshots. Will enable "Previous Versions" support for SMB. Default is true.

The `snapshot_policy` block supports:
* `daily_schedule` - (Optional) If enabled, make a snapshot every day. Defaults to midnight.
* `enabled` - (Optional) If enabled, make snapshots automatically according to the schedules. Default is false.
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


The `billing_label` block supports:
* `key` - (Required) Must be a minimum length of 1 character and a maximum length of 63 characters, and cannot be empty. Can contain only lowercase letters, numeric characters, underscores, and dashes. All characters must use UTF-8 encoding, and international characters are allowed. Must start with a lowercase letter or international character.
* `value` - (Required) Can be empty, and have a maximum length of 63 characters. Can contain only lowercase letters, numeric characters, underscores, and dashes. All characters must use UTF-8 encoding, and international characters are allowed.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the volume.

## Unique id versus name

With NetApp_GCP, every resource has a unique id. Names are not necessarily unique, but it is recommended to keep them unique per region.