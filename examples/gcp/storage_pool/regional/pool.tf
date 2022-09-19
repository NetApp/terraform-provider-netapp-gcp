# This example creates a regional Storage Pool

# local variables
locals {
  network = "cvs-prd-shared-vpc"
  pool_name = "sp_regional_example"
  region = "europe-west1"
  zone = "europe-west1-b"
  secondary_zone = "europe-west1-c"
  pool_size = 1024
  volume_name = "tf-regional-volume"
  volume_size = 100
}

resource "netapp-gcp_storage_pool" "regional-storage-pool" {
  name = local.pool_name
  region = local.region
  zone = local.zone
  secondary_zone = local.secondary_zone
  network = local.network
  size = local.pool_size
  service_level = "ZoneRedundantStandardSW"
  storage_class = "software"
}

resource "netapp-gcp_volume" "gcp-volume-nfs" {
  # This block mirrors pool settings
  region = netapp-gcp_storage_pool.regional-storage-pool.region
  zone = netapp-gcp_storage_pool.regional-storage-pool.zone
  network = netapp-gcp_storage_pool.regional-storage-pool.network
  pool_id = netapp-gcp_storage_pool.regional-storage-pool.id
  regional_ha = true
  storage_class = "software" 
  service_level = "standard"  # leave as is

  # Volume specific settings
  name = local.volume_name
  protocol_types = ["NFSv3"]
  size = local.volume_size

  # Snapshot policy definition
  snapshot_policy {
    enabled = true
  }
}