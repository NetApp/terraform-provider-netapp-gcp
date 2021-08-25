# This example shows usage of variables, outputs and tfvars files

resource "netapp-gcp_volume" "gcp-volume" {
  provider = netapp-gcp
  name = var.volume_name
  # When working with shared VPC, specify project number of host project
  # shared_vpc_project_number = "<hosting_project_number>"
  region = var.region
  protocol_types = var.protocol
  network = var.network
  size = var.size
  service_level = var.service_level
  # storage_class: choose "software for CVS, choose "hardware" for CVS-Performance
  storage_class = var.storage_class
  # zone: For storage_class = "software" specification of zone is required
  # zone = var.zone

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
  }
}
