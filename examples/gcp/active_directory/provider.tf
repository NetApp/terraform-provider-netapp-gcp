# Specify the provider and access details
provider "netapp-gcp" {
  project         = "${var.gcp_project}"
  service_account = "${var.gcp_service_account}"
}

