# This example creates an Active Directory entry for the specified region
# Please note, that CVS will not evaluate the correctness of this data
# until the first SMB volume using this entry will be created

# local variables
locals {
  region = "europe-west4"
  connection_type = "hardware"
  ad_username = "ADjoinuser"
  ad_password = "dont_show_it_to_me"
  ad_domain = "test.example.com"
  ad_dns_server = "1.2.3.4"
  ad_net_bios = "smbserver"
  ad_organizational_unit = "OU=Computers,OU=abc,OU=def,DC=sub,DC=example,DC=com"
  ad_site = ""
}

resource "netapp-gcp_active_directory" "gcp-active-directory" {
  provider = netapp-gcp
  region = local.region
  connection_type = local.connection_type
  domain = local.ad_domain
  dns_server = local.ad_dns_server
  site = local.ad_site
  organizational_unit = local.ad_organizational_unit
  net_bios = local.ad_net_bios
  aes_encryption = true

  username = local.ad_username
  password = local.ad_password
}
