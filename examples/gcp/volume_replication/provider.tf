# Specify the provider and access details
provider "netapp-gcp" {
  # numerical project number (not project ID)
  # alternatively, set GCP_PROJECT environment variable
  project         = var.gcp_project
  
  # path to JSON key file for IAM service account with "roles/netappcloudvolumes.admin" privileges
  # alternatively, set GCP_SERVICE_ACCOUNT environment variable
  service_account = var.gcp_service_account
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    netapp-gcp = {
      source = "NetApp/netapp-gcp"
      version = "~> 22.3.0"
    }
  }
}
