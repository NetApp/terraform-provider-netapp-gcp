# This example is the bare minimum to create a volume
# It creates a 1TB NFS volume in the specified region
# This can be considered as the "Hello World" example of CVS
# If it works, it means you successfully on-boarded the CVS service,
# created an service account which can talk to the CVS API and
# configured Terraform correctly to do the job

# local variables
locals {
  network = "ncv-vpc"
  volume_name = "terraform-volume-example"
  region = "europe-west3"
  size = 1024
  service_level = "standard"
}

resource "netapp-gcp_volume" "gcp-minimal-volume" {
  name = local.volume_name
  region = local.region
  protocol_types = ["NFSv3"]
  network = local.network
  # When working with shared VPC, specify project number of host project
  # shared_vpc_project_number = "<hosting_project_number>"
  size = local.size
  service_level = local.service_level
  # storage_class: choose "software for CVS, choose "hardware" for CVS-Performance
  storage_class = "hardware"
  # zone: For storage_class = "software" specification of zone is required
  # zone = "europe-west1-b"
  # when using storage_class = "software", enabling snapshot_policy is required
  # snapshot_policy {
  #   enabled = true
  # }
}
