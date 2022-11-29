# This example creates a zonal Storage Pool

# local variables
locals {
  network = "cvs-prd-shared-vpc"
  pool_name = "sp_zonal_example"
  region = "europe-west1"
  zone = "europe-west1-b"
  pool_size = 1024
  volume_name = "tf-zonal-volume"
  volume_size = 100
}

resource "netapp-gcp_storage_pool" "zonal-storage-pool" {
  name = local.pool_name
  region = local.region
  zone = local.zone
  network = local.network
  size = local.pool_size
  service_level = "StandardSW"
  storage_class = "software"
}

resource "netapp-gcp_volume" "gcp-volume-nfs" {
  # This block mirrors pool settings
  region = netapp-gcp_storage_pool.zonal-storage-pool.region
  zone = netapp-gcp_storage_pool.zonal-storage-pool.zone
  network = netapp-gcp_storage_pool.zonal-storage-pool.network
  pool_id = netapp-gcp_storage_pool.zonal-storage-pool.id
  regional_ha = false
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