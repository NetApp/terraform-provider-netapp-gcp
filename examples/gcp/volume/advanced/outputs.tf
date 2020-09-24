# returns IP/hostname of CVS export, e.g. 10.194.0.4
output "server" {
  value       = netapp-gcp_volume.gcp-volume.mount_points[0].server
  description = "The server address of the cloud volume."
}
# returns export path of CVS export, e.g. /ecstatic-dazzling-chandrasekhar
output "export" {
  value       = netapp-gcp_volume.gcp-volume.mount_points[0].export
  description = "The export path of the cloud volume."
}

output "protocol_type"  {
  value       = netapp-gcp_volume.gcp-volume.mount_points[0].protocol_type
  description = "The protocol type of the export."
}

# If volume got multiple mount points (e.g. NFSv3 and NFSv4), use below method to get data for specific protocol
output "exportfull" {
  value       = [for x in netapp-gcp_volume.gcp-volume.mount_points: "${x.server}:${x.export}" if x.protocol_type == "NFSv3"][0]
  description = "The full export path of the cloud volume."
}
