# Specify GCP resources

resource "netapp-gcp_volume" "gcp-volume" {
  provider = netapp-gcp
  name = "deleteme_asapGO"
  region = "us-west2"
  protocol_types = ["NFSv3"]
  network = "cvs-terraform-vpc"
  size = 1024
  service_level = "premium"
}

resource "netapp-gcp_snapshot" "gcp-snapshot" {
  provider = netapp-gcp
  name = "deleteme_snapshot_asapGo"
  region = "us-west2"
  volume_name =  "deleteme_asapGO"
  creation_token = "agitated-affectionate-skossi"
  depends_on = [netapp-gcp_volume.gcp-volume]
}
