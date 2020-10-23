package gcp

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGCPSnapshot_basic(t *testing.T) {

	var snapshot listSnapshotResult
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnapshotConfigCreate(VolName, Region, SnapshotName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("netapp-gcp_snapshot.gcp-snapshot-acc", &snapshot),
					resource.TestCheckResourceAttr("netapp-gcp_snapshot.gcp-snapshot-acc", "name", SnapshotName),
					resource.TestCheckResourceAttr("netapp-gcp_snapshot.gcp-snapshot-acc", "volume_name", VolName),
					resource.TestCheckResourceAttr("netapp-gcp_snapshot.gcp-snapshot-acc", "region", Region),
				),
			},
			{
				Config: testAccSnapshotConfigUpdate(VolName, Region, "update-test-snapshot"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("netapp-gcp_snapshot.gcp-snapshot-acc", &snapshot),
					resource.TestCheckResourceAttr("netapp-gcp_snapshot.gcp-snapshot-acc", "name", "update-test-snapshot"),
					resource.TestCheckResourceAttr("netapp-gcp_snapshot.gcp-snapshot-acc", "volume_name", VolName),
					resource.TestCheckResourceAttr("netapp-gcp_snapshot.gcp-snapshot-acc", "region", Region),
				),
			},
		},
	})
}

func testAccCheckSnapshotDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "netapp-gcp_snapshot" {
			continue
		}

		volume := volumeRequest{}
		volume.Region = rs.Primary.Attributes["region"]
		volume.Name = rs.Primary.Attributes["volume_name"]
		volume.CreationToken = rs.Primary.Attributes["creation_token"]

		volresult, err := client.getVolumeByNameOrCreationToken(volume)
		if err == nil {
			retriveSnapshot := listSnapshotRequest{}
			retriveSnapshot.Region = volume.Region
			retriveSnapshot.VolumeID = volresult.VolumeID
			retriveSnapshot.SnapshotID = rs.Primary.ID

			response, err := client.getSnapshotByID(retriveSnapshot)

			if err == nil {
				if response.SnapshotID != "" {
					return fmt.Errorf("Error snapshot %s still exists in %s", rs.Primary.ID, response)
				}
			}
		}
	}
	// Error: code: 500, message: Error creating volume - Cannot spawn additional jobs. Please wait for the ongoing jobs to finish and try again
	// add wait time to avoid above error after starting volume test.
	time.Sleep(30 * time.Second)
	return nil
}

// check terraform state to see if create successfully
func testAccCheckSnapshotExists(name string, snapshot *listSnapshotResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%s not found in state", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot ID is set")
		}

		volume := volumeRequest{}
		volume.Region = rs.Primary.Attributes["region"]
		volume.Name = rs.Primary.Attributes["volume_name"]
		volume.CreationToken = rs.Primary.Attributes["creation_token"]

		volresult, err := client.getVolumeByNameOrCreationToken(volume)
		if err != nil {
			return fmt.Errorf("Error getting volume ID")
		}

		retriveSnapshot := listSnapshotRequest{}
		retriveSnapshot.SnapshotID = rs.Primary.Attributes["id"]
		retriveSnapshot.Region = volume.Region
		retriveSnapshot.VolumeID = volresult.VolumeID

		var res listSnapshotResult
		res, err = client.getSnapshotByID(retriveSnapshot)
		if err != nil {
			return fmt.Errorf("Not able to get snapshot")
		}

		if res.SnapshotID != rs.Primary.ID {
			return fmt.Errorf("Snapshot id does not match")
		}

		*snapshot = res

		return nil
	}
}

const VolName = "acceptant-test-volume"
const Region = "us-east4"
const SnapshotName = "acceptant-test-snapshot"

// Create volume and snapshot based the created volume
func testAccSnapshotConfigCreate(Volume string, Location string, Snapshot string) string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "gcp-volume-acc" {
		provider = netapp-gcp
		name = "%s"
		region = "%s"
		zone = "gcp-zone"
		storage_class = "hardware"
		protocol_types = ["NFSv3"]
		network = "cvs-terraform-vpc"
		volume_path = "terraform-acceptance-test-paths"
		size = 1024
		service_level = "extreme"
	}
	
	resource "netapp-gcp_snapshot" "gcp-snapshot-acc" {
		provider = netapp-gcp
		name = "%s"
		region = "${netapp-gcp_volume.gcp-volume-acc.region}"
		volume_name =  "${netapp-gcp_volume.gcp-volume-acc.name}"
		depends_on = [netapp-gcp_volume.gcp-volume-acc] 
	}
	`, Volume, Location, Snapshot)
}

// Upate snapshot name
func testAccSnapshotConfigUpdate(Volume string, Location string, Snapshot string) string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "gcp-volume-acc" {
		provider = netapp-gcp
		name = "%s"
		region = "%s"
		zone = "gcp-zone"
		storage_class = "hardware"
		protocol_types = ["NFSv3"]
		network = "cvs-terraform-vpc"
		volume_path = "terraform-acceptance-test-paths"
		size = 1024
		service_level = "extreme"
	}
	
	resource "netapp-gcp_snapshot" "gcp-snapshot-acc" {
		provider = netapp-gcp
		name = "%s"
		region = "${netapp-gcp_volume.gcp-volume-acc.region}"
		volume_name =  "${netapp-gcp_volume.gcp-volume-acc.name}"
		depends_on = [netapp-gcp_volume.gcp-volume-acc] 
	}
	`, Volume, Location, Snapshot)
}
