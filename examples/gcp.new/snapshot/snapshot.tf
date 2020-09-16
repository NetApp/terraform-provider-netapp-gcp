# This example creates an snapshot of a volume specified via a volume data source

# local variables
locals {
    volume_name = "data-volume1"
    region = "us-central1"
    snapshot_name = "snapshot1"
}

# use data source to query an existing snapshot
data "netapp-gcp_volume" "snapshot-volume" {
    name = local.volume_name
    region = local.region
}

# create snapshot resource
resource "netapp-gcp_snapshot" "snapshot-example" {
    name = local.snapshot_name
    volume_name = local.volume_name
    region = local.region
}
