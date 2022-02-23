---
layout: "netapp_gcp"
page_title: "Provider: NetApp_GCP"
sidebar_current: "docs-netapp-gcp-index"
description: |-
  The NetApp_GCP provider is used to interact with the NetApp Cloud Volumes Service for Google Cloud resources. The provider needs to be configured with the proper credentials
  before it can be used.
---

# NetApp_GCP Provider

The NetApp_GCP provider is used to interact with the NetApp Cloud Volumes Service for Google Cloud resources.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

~> **NOTE:** The NetApp_GCP Provider currently represents _initial support_
and therefore may undergo significant changes as the community improves it.

## Example Usage
Please see our [examples on GitHub](https://github.com/NetApp/terraform-provider-netapp-gcp/tree/master/examples/gcp)

```
# Install provider from terraform register, only available for version 0.13 and above.
terraform {
  required_providers {
    netapp-gcp = {
      source = "NetApp/netapp-gcp"
      version = "20.10.0"
    }
  }
}

# Configure the NetApp_GCP Provider
provider "netapp-gcp" {
  project         = "YOUR_PROJECT_NUMBER"
  service_account = "YOUR_SERVICE_ACCOUNT"
}

# Create a volume tied to an account
resource "netapp-gcp_volume" "gcp-volume" {
  provider = netapp-gcp
  name = "deleteme_asapGO_jusitin"
  region = "us-west2"
  protocol_types = ["NFSv3"]
  network = "cvs-terraform-vpc"
  size = 1024
  service_level = "premium"
  volume_path = "deleteme-asapGO"
  snapshot_policy {
    enabled = true
    daily_schedule {
      hour = 10
      minute = 1
    }
  }
  export_policy {
    rule {
      allowed_clients = "0.0.0.0/0"
      access= "ReadWrite"
      nfsv3 {
        checked =  true
      }
      nfsv4 {
        checked = false
      }
    }
    rule {
      allowed_clients= "10.0.0.0"
      access= "ReadWrite"
      nfsv3 {
        checked =  true
      }
      nfsv4 {
        checked = false
      }
    }
  }
}

# Create a snapshot tied to an volume
resource "netapp-gcp_snapshot" "gcp-snapshot" {
  name = "deleteme_asapGO_jusitin-snapshot"
  region = "us-west2"
  volume_name =  "deleteme_asapGO_jusitin"
  # creation_token = "unique-token-number"
}

# Create an Active Directory
resource "netapp-gcp_active_directory" "gcp-active-directory" {
  provider = netapp-gcp
  region = "us-west2"
	username = "test_user"
	password = "netapp"
  domain = "example.com"
  dns_server = "10.0.0.0"
  net_bios = "cvsserver"
}

# Create a volume_backup tied to an volume
resource "netapp-gcp_volume_backup" "gcp-volume-backup" {
  name = "main-volume-backup"
  region = "us-west2"
  volume_name =  "main-volume"
  creation_token = "unique-token-number"
}
```

## Argument Reference

The following arguments are used to configure the NetApp_GCP Provider:

* `project` - (Required) This is the project number or project ID for NetApp_GCP API operations.
* There are three ways for the NetApp_GCP API operation authentication. One of them should be used in the provider.
* `credentials` - (Optional) JSON key as base64-encoded string of the account for NetApp_GCP API operations.
* `service_account` - (Optional) There are two ways to be used:
  - Using the service account key file for the authentication. This is the JSON file path of the service_account with the "roles/netappcloudvolumes.admin" privileges.
  - Using the service account impersonation for the authentication. This is the service account email address service-account-name@the-project-id.iam.gserviceaccount.com using the impersonation for NetApp_GCP API operations. The account running the terraform needs be attached to this service account with "role/serviceAccountTokenCreator". The service account needs to have "roles/netappcloudvolumes.admin" privileges.
* `token_duration` - (Optional) The token life duration in minutes in the service account impersonation case. Default is 60.

## Required Privileges

These settings were tested with GCP Google Cloud SDK 274.0.0.
For additional information on roles and permissions, please refer to official
NetApp_GCP documentation.

