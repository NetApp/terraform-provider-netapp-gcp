## 22.6.1
BUG FIXES:

* resource/volume: remove check for regionalHA and zone when storageClass is software.
* resource/volume: update the create and delete volume error message verification.

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
