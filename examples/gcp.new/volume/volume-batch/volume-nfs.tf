# This example creates multiple volumes from a CSV file
# If you run into errors saying
# "Error: code: 500, message: Error creating volume - Cannot spawn additional jobs. Please wait for the ongoing jobs to finish and try again"
# you are seeing Issue #13 (https://github.com/NetApp/terraform-provider-netapp-gcp/issues/13)
# until it is fixed, use "terraform apply --parallelism=1"

# local variables

locals {
  # read CVS file
  volumes = csvdecode(file("${path.module}/volumes.csv"))
  network = "ncv-vpc"
}

# Specify GCP resources
resource "netapp-gcp_volume" "gcp-volumes-batch" {
  count = length(local.volumes)

  name = local.volumes[count.index].volume_name
  volume_path = local.volumes[count.index].volume_name
  region = local.volumes[count.index].region
  protocol_types = ["NFSv3"]
  network = local.network
  size = local.volumes[count.index].size
  service_level = local.volumes[count.index].service_level
  
  snapshot_policy {
    enabled = true
    daily_schedule {
      snapshots_to_keep = 30
      hour = 1
      minute = 0
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
