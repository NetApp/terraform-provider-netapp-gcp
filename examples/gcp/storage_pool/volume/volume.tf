# This example creates a volume using an existing zonal Storage Pool
# The network, region and zone arguments need to match with the Storage Pool values

# local variables
locals {
  volume_name = "terraform-volume-example"
  size = 50
  pool_id = "12ab34cd-56ef-78gh-90ij-12kl34mn56op"
  network = "ncv-vpc"
  region = "europe-west1"
  zone = "europe-west1-b"
}

resource "netapp-gcp_volume" "cvs-volume" {
  name = local.volume_name
  size = local.size
  pool_id = local.pool_id
  protocol_types = ["NFSv3"]
  storage_class = "software"
}
