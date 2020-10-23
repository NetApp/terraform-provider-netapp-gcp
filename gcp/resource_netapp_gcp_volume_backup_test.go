package gcp

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVolumeBackup_basic(t *testing.T) {

	var volume listVolumeBackupResult
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCPVolumeBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeBackupConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPVolumeBackupExists("netapp-gcp_volume_backup.gcp-volume-backup", &volume),
					testCheckResourceAttr("netapp-gcp_volume_backup.gcp-volume-backup", "name", "terraform-acceptance-test-1"),
					testCheckResourceAttr("netapp-gcp_volume_backup.gcp-volume-backup", "region", "us-east4"),
					testCheckResourceAttr("netapp-gcp_volume_backup.gcp-volume-backup", "volume_name", "test-volume"),
				),
			},
		},
	})
}

func testAccCheckGCPVolumeBackupDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	var volumeID string
	volumeBackup := listVolumeBackupRequest{}

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "netapp-gcp_volume" {
			volumeID = rs.Primary.ID
			response, err := client.getVolumeByID(volumeRequest{
				VolumeID: volumeID,
				Region:   rs.Primary.Attributes["region"],
			})
			if err == nil {
				if response.LifeCycleState != "deleted" {
					return fmt.Errorf("volume (%s) still exists", response.VolumeID)
				}
			}
		} else if rs.Type == "netapp-gcp_volume_backup" {
			volumeBackup.Region = rs.Primary.Attributes["region"]
			volumeBackup.VolumeBackupID = rs.Primary.ID
		}
	}
	volumeBackup.VolumeID = volumeID
	var response listVolumeBackupResult
	response, err := client.getVolumeBackupByID(volumeBackup)
	if err != nil {
		return err
	}
	if err == nil {
		if response.VolumeBackupID != "" {
			return fmt.Errorf("volume (%s) still exists", response.VolumeBackupID)
		}
	}
	return nil
}

func testAccCheckGCPVolumeBackupExists(name string, volumeBU *listVolumeBackupResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}

		volumeBackup := listVolumeBackupRequest{}

		volumeBackup.Region = rs.Primary.Attributes["region"]

		volume := volumeRequest{}
		volume.Region = volumeBackup.Region
		volume.Name = rs.Primary.Attributes["volume_name"]
		volume.CreationToken = rs.Primary.Attributes["creation_token"]

		volresult, err := client.getVolumeByNameOrCreationToken(volume)
		if err != nil {
			log.Print("Error getting volume ID")
			return err
		}

		volumeBackup.VolumeID = volresult.VolumeID

		volumeBackup.VolumeBackupID = rs.Primary.ID
		var response listVolumeBackupResult
		response, err = client.getVolumeBackupByID(volumeBackup)
		if err != nil {
			return err
		}

		if response.VolumeBackupID != rs.Primary.ID {
			return fmt.Errorf("Resource ID and volume ID do not match")
		}

		*volumeBU = response

		return nil
	}
}

func testAccVolumeBackupConfigCreate() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "gcp-volume-acc" {
		provider = netapp-gcp
		name = "test-volume"
		region = "us-east4"
		zone = "gcp-zone"
		storage_class = "hardware"
		protocol_types = ["NFSv3"]
		network = "cvs-terraform-vpc"
		volume_path = "terraform-acceptance-test-paths"
		size = 1024
		service_level = "extreme"
	}

	resource "netapp-gcp_volume_backup" "gcp-volume-backup" {
		name = "terraform-acceptance-test-1"
		region = "us-east4"
		volume_name = "${netapp-gcp_volume.gcp-volume-acc.name}"
		creation_token = "${netapp-gcp_volume.gcp-volume-acc.volume_path}"
		depends_on = [netapp-gcp_volume.gcp-volume-acc]
	  }
  `)
}
