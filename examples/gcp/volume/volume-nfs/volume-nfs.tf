# This example creates an NFS volume with export policy and automatic scnapshots schedule

# local variables
locals {
  network = "ncv-vpc"
  volume_name = "terraform-nfs-example"
  region = "europe-west3"
  size = 1024
  service_level = "premium"
}

resource "netapp-gcp_volume" "gcp-volume-nfs" {
  name = local.volume_name
  region = local.region
  protocol_types = ["NFSv3"]
  network = local.network
  size = local.size
  service_level = local.service_level
 
  # up to 5 export rules
  export_policy {
    rule {
      allowed_clients = "0.0.0.0/0"
      access= "ReadOnly"
      nfsv3 {
        checked =  true
      }
      nfsv4 {
        checked = false
      }
    }
    rule {
      allowed_clients = "10.10.20.0/24"
      access= "ReadWrite"
      nfsv3 {
        checked =  true
      }
      nfsv4 {
        checked = false
      }
    }
  }

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

