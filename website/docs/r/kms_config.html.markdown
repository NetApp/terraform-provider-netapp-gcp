---
layout: "netapp_gcp"
page_title: "NetApp_GCP: netapp_gcp_kms_config"
sidebar_current: "docs-netapp-gcp-resource-kms-config"
description: |-
  Provides a NetApp_GCP kms configuration resource. This can be used to create a new kms config on the GCP-CVS.
---

# netapp_gcp\_kms

Provides a NetApp_GCP kms configuration resource. Volumes are encrypted at rest. The type of key used is determined at volume creation. 

To use Customer-managed Encryption Keys (CMEK), use this resource to use a CMEK from Google Cloud KMS instead. If no CMEK definition exists for the region, a service-managed key will be used instead.

## Example Usages

**Create NetApp_GCP kms config:**

```
resource "netapp-gcp_kms_config" "kms-example" {
    provider = netapp-gcp
    region = "us-central1"
    key_project_id = "central-kms-key-project"
    key_ring_location = "global"
    key_ring_name = "TheOneRing"
    key_name = "VirginiaKey"
    network = "projects/{project_number}/global/networks/{vpc_name}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) Name of the region to create a KMS config for.
* `key_name` - (Required, modifiable) Name of the key to be used for encryption. This key should be in the keyRing mentioned in keyRing field.
* `key_project_id` - (Optional,modifiable) Project ID of project where the key to be used for encryption is residing. Use if key is located in different project.
* `key_ring_location` - (Required) Location/region of the keyRing.
* `key_ring_name` - (Required, modifiable) Key ring containing the keys to be used for volume encryption.
* `network` - (Required) The path of the VPC the volumes are attached to, for example: "projects/12345678/global/networks/my-shared-vpc".

  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the kms.