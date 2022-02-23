# Specify the provider and access details
provider "netapp-gcp" {
  #
  # numerical project number or project ID xxx-xxx-xxx
  # alternatively, set GCP_PROJECT environment variable
  # project         = "123456890"
  #
  #
  # JSON key as base64-encoded string
  # alternatively, set GCP_CREDENTIALS environment variable
  # credentials = "xxxxxxxxxxxxxxxxxxxxxxx"
  #
  #
  # service_account
  # alternatively, set GCP_SERVICE_ACCOUNT environment variable
  #
  #
  # When using a service account key for the authorization
  # Use path to JSON key file for IAM service account with "roles/netappcloudvolumes.admin" privileges
  # service_account = "/Users/abc/key.json"
  #
  #
  # When using a service account impersonation for the authentication:
  # service_account is service account email xxxx@<project id>.iam.gserviceaccount.com
  # service_account = "serviceaccount@the-project-id.iam.gserviceaccount.com"
  #
  # token life duration in minutes in the service_account impersonation case
  # alternatively, set GCP_TOKEN_DURATION environment variable
  # token_duration = 1
  #
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
