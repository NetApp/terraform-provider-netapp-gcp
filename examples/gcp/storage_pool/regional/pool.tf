# This example creates a regional Storage Pool

# local variables
locals {
  network = "ncv-vpc"
  pool_name = "sp_regional_example"
  region = "europe-west1"
  zone = "europe-west1-b"
  secondary_zone = "europe-west1-c"
  size = 1024
}

resource "netapp-gcp_storage_pool" "regional-storage-pool" {
  name = local.pool_name
  region = local.region
  zone = local.zone
  secondary_zone = local.secondary_zone
  network = local.network
  size = local.size
  service_level = "ZoneRedundantStandardSW"
  storage_class = "software"
}