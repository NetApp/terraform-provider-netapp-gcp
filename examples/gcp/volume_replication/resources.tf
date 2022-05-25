# This example creates two NFS volumes in different regions and
# sets up a Volume Replication

# local variables
locals {
  source_volume_name = "terraform-nfs-source"
  source_region = "europe-west3"
  source_size = 1024
  source_service_level = "premium"

  destination_volume_name = "terraform-nfs-destination"
  destination_region = "europe-west4"
  destination_size = 1024
  destination_service_level = "standard"
}

# Set up source volume
resource "netapp-gcp_volume" "gcp-volume-source" {
  name = local.source_volume_name
  region = local.source_region
  protocol_types = ["NFSv3"]
  network = var.network
  size = local.source_size
  service_level = local.source_service_level
  type_dp = "false"
}

# Set up secondary volume (= destination)
resource "netapp-gcp_volume" "gcp-volume-destination" {
  name = local.destination_volume_name
  region = local.destination_region
  protocol_types = ["NFSv3"]
  network = var.network
  size = local.destination_size
  service_level = local.destination_service_level
  type_dp = "true"
}

# Setting up the replication
resource "netapp-gcp_volume_replication" "gcp-volume-replication" {
  provider = netapp-gcp
  name = "myReplica"
  endpoint_type = "dst"
  # if endpoint_type="dst"
  #  region=<region_of_destination_volume>
  #  remote_region=<region_of_source_volume>
  # if endpoint_type="src" (Note: not supported currently)
  #  region=<region_of_source_volume>
  #  remote_region=<region_of_destination_volume>
  region = local.destination_region
  remote_region = local.source_region

  source_volume_id = netapp-gcp_volume.gcp-volume-source.id
  destination_volume_id = netapp-gcp_volume.gcp-volume-destination.id
  schedule = "hourly"
}
