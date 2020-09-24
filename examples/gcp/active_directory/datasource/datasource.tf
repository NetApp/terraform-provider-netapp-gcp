# This example shows how to query for an Active Directory entry
# Use case: Only allow creation of SMB volume, if AD entry exists

locals {
  region = "europe-west3"
}

data "netapp-gcp_active_directory" "ad-europe-west3" {
    # domain = "example.com"
    region = local.region
}

# Some example outputs
output "domain_name" {
  value       = data.netapp-gcp_active_directory.ad-europe-west3.domain
  description = "Print name of AD domain"
}

output "domain_region" {
  value       = data.netapp-gcp_active_directory.ad-europe-west3.region
  description = "Print region of AD domain"
}

output "domain_organizational_unit" {
  value       = data.netapp-gcp_active_directory.ad-europe-west3.organizational_unit
  description = "Print Organization Unit (OU) of AD domain"
}
