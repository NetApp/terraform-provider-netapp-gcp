# This example creates a zonal Storage Pool

# local variables
locals {
  network = "ncv-vpc"
  pool_name = "sp_zonal_example"
  region = "europe-west1"
  zone = "europe-west1-b"
  size = 1024
}

resource "netapp-gcp_storage_pool" "zonal-storage-pool" {
  name = local.pool_name
  region = local.region
  zone = local.zone
  network = local.network
  size = local.size
  service_level = "StandardSW"
  storage_class = "software"
}