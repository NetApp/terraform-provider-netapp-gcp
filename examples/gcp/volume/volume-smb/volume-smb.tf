# This example creates an SMB volume with automatic scnapshots schedule

# local variables
locals {
  network = "ncv-vpc"
  volume_name = "terraform-smb-example"
  region = "europe-west3"
  size = 2048
  service_level = "extreme"
}

data "netapp-gcp_active_directory" "myad" {
    region = local.region
}

resource "netapp-gcp_volume" "gcp-volume-smb" {
  name = local.volume_name
  region = local.region
  protocol_types = ["SMB"]
  network = local.network
  size = local.size
  service_level = local.service_level

  # Advice: Since SMB volumes can only be created if an Active Directory connection exists for the region,
  # depend the SMB volume on the AD resource. Either create the AD from TF, or use AD data source to query for
  # existing AD connection
  depends_on = [
    data.netapp-gcp_active_directory.myad
  ]

  # Snapshot policy definition
  snapshot_policy {
    enabled = true
    hourly_schedule {
      snapshots_to_keep = 48
      minute = 1
    }
    daily_schedule {
      snapshots_to_keep = 14
      hour = 23
      minute = 2
    }
    weekly_schedule {
      snapshots_to_keep = 4
      hour = 1
      minute = 3
      day = "Monday"
    }
    monthly_schedule {
      snapshots_to_keep = 6
      hour = 2
      minute = 4
      days_of_month = 6
    }    
  }
}
