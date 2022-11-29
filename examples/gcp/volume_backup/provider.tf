# Specify the provider and access details
provider "netapp-gcp" {
  # numerical project number (not project ID)
  # alternatively, set GCP_PROJECT environment variable
  # project         = "123456890"
  
  # path to JSON key file for IAM service account with "roles/netappcloudvolumes.admin" privileges
  # alternatively, set GCP_SERVICE_ACCOUNT environment variable
  # service_account = "/Users/abc/key.json"
}

provider "netapp-gcp" {
  project         = "775487933698"
  service_account = "/Users/swenjun/Downloads/cvs-terraform-dev-26fd2c8c8343.json"
}
