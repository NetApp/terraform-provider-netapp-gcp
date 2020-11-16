variable "gcp_project" {
  type = string
  description = "Project Number."
  default = null
}

variable "gcp_service_account" {
  type = string
  description = "Absolute file path to service account .json key file."
  default = null
}

variable "network" {
  description = "Network to deploy to. Only one of network or subnetwork should be specified."
  default     = ""
}

variable "region" {
  type        = string
  description = "Region where the instances should be created."
  default     = null
}

variable "volume_name" {
  type        = string
  description = "Name of CVS volume."
  default     = null
}

variable "protocol" {
  type        = list(string)
  description = "Enabled NAS protocols NFSv3, NFSv4, CIFS, SMB."
  default     = ["NFSv3"]
}

variable "size" {
  type        = number
  description = "Size of volume in GB"
  default     = 1024
}

variable "service_level" {
  type        = string
  description = "Service level low, medium or high."
  default     = "medium"
}

variable "storage_class" {
  type        = string
  description = "Type of CVS service: CVS=software, CVS-Performance=hardware."
  default     = "hardware"
}

variable "zone" {
  type        = string
  description = "GCP zone CVS-Software is deployed to. Required for CVS-Software."
  default     = null
}

variable "ad_username" {
  type        = string
  description = "Active Directory username for joining domain."
  default     = "Administrator"
}

variable "ad_password" {
  type        = string
  description = "Password for Active Directory username for joining domain."
  default     = null
}

variable "ad_domain" {
  type        = string
  description = "Name of Active Directory domain."
  default     = null
}

variable "ad_dns" {
  type        = string
  description = "IP of Active Directory DNS server."
  default     = null
}

variable "ad_netbios" {
  type        = string
  description = "Netbios name for SMB server."
  default     = null
}

variable "ad_organizational_unit" {
  type        = string
  description = "Organizational Unit."
  default     = "CN=Computers"
}

variable "ad_site" {
  type        = string
  description = "Active Directory Site"
  default     = ""
}
