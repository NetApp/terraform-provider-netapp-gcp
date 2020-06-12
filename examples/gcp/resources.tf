# Specify GCP resources

data "netapp-gcp_volume" "data-volume" {
  name = "vol"
  region = "us-east4"
}

output "volume" {
  value = data.netapp-gcp_volume.data-volume
}

resource "netapp-gcp_volume" "gcp-volume" {
  provider = netapp-gcp
  name = "deleteme_asapGO"
  region = "us-east4"
  protocol_types = ["NFSv3"]
  network = "cvs-terraform-vpc"
  size = 1024
  volume_path = "deleteme-asapGO"
  service_level = "standard"
}

resource "netapp-gcp_snapshot" "gcp-snapshot" {
  provider = netapp-gcp
  name = "deleteme_snapshot_asapGo"
  region = "us-east4"
  volume_name =  data.netapp-gcp_volume.data-volume.name
  creation_token = data.netapp-gcp_volume.data-volume.volume_path
  depends_on = [netapp-gcp_volume.gcp-volume]
}
