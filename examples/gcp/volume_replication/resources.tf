
terraform {
  required_version = ">= 0.13"
  required_providers {
   netapp-gcp = {
   source = "github.com/netapp/netapp-gcp"
   version = "1.0.0"
   }
  }
 }

# Configure the NetApp_GCP Provider
provider "netapp-gcp" {
  project         = "project_id"
  service_account = "/Path-to-file/cvs-terraform.json"
}


# Specify GCP resources

resource "netapp-gcp_volume_replication" "gcp-volume-replication" {
  provider = netapp-gcp
  name = "myReplica"
  region = "us-east4"
  remote_region = "us-west2"
  destination_volume_id = "dst-volume-id"
  source_volume_id = "src-volume-id"
  endpoint_type = "dst"
  schedule = "hourly"
}
