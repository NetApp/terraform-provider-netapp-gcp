# Specify the provider and access details
provider "netapp-gcp" {
  # numerical project number or project ID xxx-xxx-xxx
  # alternatively, set GCP_PROJECT environment variable
  # project         = "123456890"
  
  # path to JSON key file for IAM service account with "roles/netappcloudvolumes.admin" privileges
  # alternatively, set GCP_SERVICE_ACCOUNT environment variable
  # service_account = "/Users/abc/key.json"
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    netapp-gcp = {
      source = "netapp.com/netapp/netapp-gcp"
      version = "~> 20.10.0"
    }
  }
}
