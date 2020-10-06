# This example creates a volume backup

# local variables
locals {
  backup_volume_name = "main-volume-backup"
  volume_name = "main-volume"
  region = "us-east4"
  creation_token = "quirky-musing-swartz"
}

resource "netapp-gcp_volume_backup" "gcp-volume-backup" {
  name = local.backup_volume_name
  region = local.region
  volume_name =  local.volume_name
  creation_token = local.creation_token
}
