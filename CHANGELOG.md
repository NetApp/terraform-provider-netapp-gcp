## 24.6.0
BUG FIXES:
* resource/volume_replication: fix update.

## 24.4.0
BUG FIXES:
* resource/volume: ([#90](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/90)), ([#92](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/92))

## 23.11.0
BUG FIXES:
* resource/active_directory: support terraform import.
* resorce/volume: fix root access issue.([#90](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/90))

## 23.5.1
ENHANCEMENTS:
* add `security_style` option in volume.html.markdown.

## 23.4.0
* resource/storage_pool: Support terraform import with "poolID:region". ([#88](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/88))

## 23.1.0
* resource/storage_pool: update and add acceptance test. ([#85](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/85))
* resource/volume: Fix for API retry timeout. ([#86](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/86))

## 22.12.0
* data_source/volume: update to reflect the change of volume resource. ([#83](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/83))
* resource/volume: add `security_style` option. ([#83](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/83))
* resource/volume: add `snapshot_id` option. 

## 22.10.0
* resource/active_directory: ([#76](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/76))

## 22.8.1
BUG FIXES:
* resource/storage_pool: Fix creation error with shared vpc. ([#69](https://github.com/NetApp/terraform-provider-netapp-gcp/issues/69))

## 22.8.0
BUG FIXES:
* resource/volume: update the create and delete volume error message verification.
* resource/volume: Updated Volume size from 1 Gib. ([#67](https://github.com/NetApp/terraform-provider-netapp-gcp/pull/67))

ENHANCEMENTS:
* Improvement on documentations. ([#68](https://github.com/NetApp/terraform-provider-netapp-gcp/pull/68))

## 22.6.1
BUG FIXES:

* resource/volume: remove check for regionalHA and zone when storageClass is software.


## 22.6.0
ENHANCEMENTS:

* resource/active_directory: add connection_type and ad_server options.

## 22.4.0
ENHANCEMENTS:

* resource/volume: remove `snap_reserve` option.
* resource/volume: without either enable NFSv3 or NFSv4, the export rule is invalid.

## 22.3.0
NEW FEATURES:

* resource/kms: create,update and delete kms config.

## 22.2.0
ENHANCEMENTS:

* resource/volume: add pool_id option.
* Support service account principal name when using service account impersonation.

BUG FIXES:

* Fix use default credentials when providing project ID. 

## 22.1.1
ENHANCEMENTS:

* resource/volume: add billing_label option.

## 20.10.0 (Oct 2020)

* **New DataSource:** netapp-gcp_active_directory
* **New Resource:** `netapp-gcp_volume_backup`
* **Updated Resource:** `netapp-gcp_volume` to support `type_dp`
* **Updated Resource:** `netapp-gcp_volume` to support `zone` and `storage_class` for SDS

## 0.1.1 (Aug 12, 2020)

* Released on Terraform Registry in addition to GitHub

* **New DataSource:** netapp-gcp_volume

## 0.1.0 (Mar 25, 2020)

FEATURES:

* **New Resource:** `netapp-gcp_volume`
* **New Resource:** `netapp-gcp_snapshot`
* **New Resource:** `netapp-gcp_active_directory`
