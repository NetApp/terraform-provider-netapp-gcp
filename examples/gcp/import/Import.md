# How to use Terraform import

The provider supports importing existing resource to put them under Terraform management. 

Importing works like this:

1. Create an empty resource file, e.g.
```tf
resource "netapp-gcp_volume" "myvolume" {
}

```
2. Find the ID of the resource to import (e.g. in UI)
3. Run the import
```bash
terraform import netapp-gcp_volume.myvolume 1bc88bc6-cc7d-5fe3-3737-8e635fe2f996
```
Since IDs are region specific, there is a very minimal chance that two volumes with the same ID exist. If the ID is no unique, the import will fail and propose another import syntax to use ID:region instead (e.g. 1bc88bc6-cc7d-5fe3-3737-8e635fe2f996:us-central1).

4. Export the output of `terraform show` into a new tf file and remove the `id` line. E.g:
```bash
terraform show -no-color | sed '/^[[:blank:]]*id /d' > myvolume.tf
```
5. Use `terraform plan` to verify the resources file against the deployed volume. It should report no changes.