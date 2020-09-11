# NetApp NetApp_GCP 0.1.2 Example

This repository is designed to demonstrate the capabilities of the [Terraform
NetApp NetApp_GCP Provider][ref-tf-netapp-gcp] at the time of the 0.1.2 release.

[ref-tf-netapp_gcp]: https://www.terraform.io/docs/providers/netapp/gcp/index.html

This example performs the following:

* Creates a number of volumes on the GCP CVS tied to the project,
  using the [`netapp_gcp_volume` resource][ref-tf-netapp-gcp-volume].

[ref-tf-netapp-gcp-volume]: https://www.terraform.io/docs/providers/netapp/gcp/r/volume.html

## Requirements

* A working GCP NetApp account.

## Usage Details

You can either clone the entire
[terraform-provider-netapp_gcp][ref-tf-netapp-gcp-github] repository, or download the
`provider.tf`, `variables.tf`, `resources.tf`, and
`terraform.tfvars.example` files into a directory of your choice. Once done,
edit the `terraform.tfvars.example` file, populating the fields with the
relevant values, and then rename it to `terraform.tfvars`. Don't forget to
configure your endpoint and credentials by either adding them to the
`provider.tf` file, or by using enviornment variables. See
[here][ref-tf-netapp-gcp-provider-settings] for a reference on provider-level
configuration values.

[ref-tf-netapp-gcp-github]: https://github.com/terraform-providers/terraform-provider-netapp-gcp
[ref-tf-netapp-gcp-provider-settings]: https://www.terraform.io/docs/providers/netapp/gcp/index.html#argument-reference

Once done, run `terraform init`, and `terraform plan` to review the plan, then
`terraform apply` to execute. If you use Terraform 0.11.0 or higher, you can
skip `terraform plan` as `terraform apply` will now perform the plan for you and
ask you confirm the changes.
