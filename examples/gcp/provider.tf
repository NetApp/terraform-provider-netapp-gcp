# Specify the provider and access details
provider "netapp-gcp" {
  # specify projectID or project number
  # alternatively, set GCP_PROJECT environment variable
  # project         = "123456890"

  # We support two ways to authenticate. Choose one or the other

  # Auth option 1: using service account key file
  # path to JSON key file for IAM service account with "roles/netappcloudvolumes.admin" privileges
  # alternatively, set GCP_SERVICE_ACCOUNT environment variable

  # service_account = "/Users/abc/key.json"

  # Auth option 2: provide credentials as text string
  # Basically providing the content of the keyfile of option 1 via the credentials
  # attribute. This is often used in Terraform Enterprise environments

  # credentials = "${base64decode(data.vault_generic_secret.gcp_key.data["private_key"])}"
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    netapp-gcp = {
      source = "NetApp/netapp-gcp"
      version = "~> 21.2.0"
    }
  }
}
