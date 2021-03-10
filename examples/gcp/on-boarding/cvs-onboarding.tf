# This code is provided as-is an an example on how CVS can be on-boarded.
# It is not supported by NetApp, since on-boarding is done via GCP APIs.
# It uses the official Googles Terraform provider to do all the work.
# See https://registry.terraform.io/providers/hashicorp/google/latest/docs for documentation

locals {
    project = "myproject-1234"
    network = "default"
    address_name = "netapp-addresses-${local.network}"

    address_ip = "192.168.200.0"   # RFC1918
    address_prefix = "24"
}

# Enable APIs
# Unfortunately, this works only if service is already enabled. If it isn't enabled, this times out
# cloudvolumesgcp-api.netapp.com needs to be enabled via Googles Marketplace
# See https://console.cloud.google.com/marketplace/product/endpoints/cloudvolumesgcp-api.netapp.com
# the other two APIs can be enabled with "gcloud services enable <api_name>"
resource "google_project_service" "gcp_apis" {
    for_each = toset([
        "servicenetworking.googleapis.com",
        "servicemanagement.googleapis.com",
        "cloudvolumesgcp-api.netapp.com",
    ])
  service = each.key

  project = local.project
  disable_dependent_services = true
  disable_on_destroy = false
}

# Create global compute address reservation for CVS to use
# gcloud compute addresses create <...> --global --addresses <...> --purpose=VPC_PEERING --prefix-length=<...> --network=<...> --no-user-output-enabled
resource "google_compute_global_address" "cvs_address_pool" {
    project = local.project
    name = local.address_name
    address = local.address_ip
    prefix_length = local.address_prefix
    ip_version = "IPV4"
    address_type = "INTERNAL"
    purpose = "VPC_PEERING"
    network = local.network
}

data "google_compute_network" "myvpc" {
    project = local.project
    name = local.network
}

# gcloud services vpc-peerings connect --service=cloudvolumesgcp-api-network.netapp.com --ranges=<...> --network=<...> --no-user-output-enabled
resource "google_service_networking_connection" "cvs_performance_peering" {
    network                 = data.google_compute_network.myvpc.self_link
    service                 = "cloudvolumesgcp-api-network.netapp.com"
    reserved_peering_ranges = [google_compute_global_address.cvs_address_pool.name]
}

data "google_compute_network" "remotevpc" {
    project = local.project
    name = local.network
}

# gcloud compute networks peerings update <...> --network=<...> --import-custom-routes --export-custom-routes
resource "google_compute_network_peering_routes_config" "cvs_routes_update" {
    project = local.project
    peering = google_service_networking_connection.cvs_performance_peering.peering
    network = local.network

    import_custom_routes = true
    export_custom_routes = true
}
