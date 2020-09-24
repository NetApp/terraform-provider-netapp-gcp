# This example creates two NFS volumes in different regions and
# sets up a Volume Replication

# local variables
locals {
  network = "ncv-vpc"

  source_volume_name = "terraform-nfs-source"
  source_region = "europe-west3"
  source_size = 1024
  source_service_level = "premium"

  destination_volume_name = "terraform-nfs-destination"
  destination_region = "europe-west2"
  destination_size = 1024
  destination_service_level = "standard"
}

# Set up source volume
resource "netapp-gcp_volume" "gcp-volume-source" {
  name = local.source_volume_name
  region = local.source_region
  protocol_types = ["NFSv3"]
  network = local.network
  size = local.source_size
  service_level = local.source_service_level
  type_dp = "false"
}

# Set up secondary volume (= destination)
resource "netapp-gcp_volume" "gcp-volume-destination" {
  name = local.destination_volume_name
  region = local.destination_region
  protocol_types = ["NFSv3"]
  network = local.network
  size = local.destination_size
  service_level = local.destination_service_level
  type_dp = "true"
}

# Setting up the replication will follow as soon it is in the provider
