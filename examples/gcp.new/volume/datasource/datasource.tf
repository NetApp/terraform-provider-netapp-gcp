# This example demonstrates how to query an existing volume and use the results
# Possible uses are creating/deleting a snapshot
# In this example, mount paths are displayed as outputs
# A more advance example would be to use the path to mount the volume to a VM created by Terraform

data "netapp-gcp_volume" "myvolume" {
    name = "data-volume1"
    region = "us-central1"
}

# returns IP/hostname of CVS export, e.g. 10.194.0.4
output "server" {
  value       = data.netapp-gcp_volume.myvolume.mount_points[0].server
  description = "The server address of the cloud volume."
}
# returns export path of CVS export, e.g. /ecstatic-dazzling-chandrasekhar
output "export" {
  value       = data.netapp-gcp_volume.myvolume.mount_points[0].export
  description = "The export path of the cloud volume."
}

output "protocol_type"  {
  value       = data.netapp-gcp_volume.myvolume.mount_points[0].protocol_type
  description = "The protocol type of the export."
}

# If volume got multiple mount points (e.g. NFSv3 and NFSv4), use below method to get data for specific protocol
output "exportfull" {
  value       = [for x in data.netapp-gcp_volume.myvolume.mount_points: "${x.server}:${x.export}" if x.protocol_type == "NFSv3"][0]
  description = "The full export path of the cloud volume."
}