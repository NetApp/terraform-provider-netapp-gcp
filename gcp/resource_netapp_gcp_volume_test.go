package gcp

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVolume_basic(t *testing.T) {

	var volume volumeResult
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCPVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPVolumeExists("netapp-gcp_volume.terraform-acceptance-test-1", &volume),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "name", "terraform-acceptance-test-1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "size", "1024"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "region", "us-west2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "service_level", "extreme"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "volume_path", "terraform-acceptance-test-path"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "protocol_types.0", "NFSv3"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.hourly_schedule.0.snapshots_to_keep", "48"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.hourly_schedule.0.minute", "1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.snapshots_to_keep", "14"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.minute", "2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.hour", "23"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.snapshots_to_keep", "4"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.minute", "3"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.hour", "1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.day", "Monday"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.snapshots_to_keep", "6"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.minute", "4"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.hour", "2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.days_of_month", "6"),
				),
			},
			{
				Config: testAccVolumeConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPVolumeExists("netapp-gcp_volume.terraform-acceptance-test-1", &volume),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "name", "terraform-acceptance-test-1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "size", "2048"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "region", "us-west2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "service_level", "extreme"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.hourly_schedule.0.snapshots_to_keep", "9"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.hourly_schedule.0.minute", "2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.snapshots_to_keep", "20"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.minute", "10"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.hour", "22"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.snapshots_to_keep", "15"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.minute", "35"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.hour", "2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.weekly_schedule.0.day", "Tuesday"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.snapshots_to_keep", "10"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.minute", "5"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.hour", "3"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.monthly_schedule.0.days_of_month", "15"),
				),
			},
			// remove temporarily since us-west2 is not working.
			// {
			// 	Config: testAccVolumeConfigCreateSMB(),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckGCPVolumeExists("netapp-gcp_volume.terraform-acceptance-test-1-SMB", &volume),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "name", "terraform-acceptance-test-1-SMB"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "size", "1024"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "region", "us-east4"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "service_level", "extreme"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "protocol_types.0", "SMB"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "snapshot_policy.0.daily_schedule.0.hour", "10"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "snapshot_policy.0.daily_schedule.0.minute", "1"),
			// 		testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1-SMB", "snapshot_policy.0.daily_schedule.0.snapshots_to_keep", "0"),
			// 	),
			// },
		},
	})
}

func testAccCheckGCPVolumeDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "netapp-gcp_volume" {
			continue
		}
		response, err := client.getVolumeByID(volumeRequest{
			VolumeID: rs.Primary.ID,
			Region:   rs.Primary.Attributes["region"],
		})
		if err == nil {
			if response.VolumeID != "" {
				return fmt.Errorf("volume (%s) still exists.", response.VolumeID)
			}

		}
	}
	return nil
}

func testAccCheckGCPVolumeExists(name string, volume *volumeResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		response, err := client.getVolumeByID(volumeRequest{
			VolumeID: rs.Primary.ID,
			Region:   rs.Primary.Attributes["region"],
		})

		if err != nil {
			return err
		}

		if response.VolumeID != rs.Primary.ID {
			return fmt.Errorf("Resource ID and volume ID do not match")
		}

		*volume = response

		return nil
	}
}

func testAccVolumeConfigCreate() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "terraform-acceptance-test-1" {
		provider = netapp-gcp
		name = "terraform-acceptance-test-1"
		region = "us-west2"
		protocol_types = ["NFSv3"]
		network = "cvs-terraform-vpc"
		size = 1024
		service_level = "extreme"
		volume_path = "terraform-acceptance-test-path"
		snapshot_policy {
			enabled = true
			hourly_schedule {
			  snapshots_to_keep = 48
			  minute = 1
			}
			daily_schedule {
			  snapshots_to_keep = 14
			  hour = 23
			  minute = 2
			}
			weekly_schedule {
			  snapshots_to_keep = 4
			  hour = 1
			  minute = 3
			  day = "Monday"
			}
			monthly_schedule {
			  snapshots_to_keep = 6
			  hour = 2
			  minute = 4
			  days_of_month = 6
			}    
		}
		export_policy {
		  rule {
			allowed_clients = "0.0.0.0/0"
			access= "ReadWrite"
			nfsv3 {
			   checked =  true
			}
			nfsv4 {
			   checked = false
			}
		  }
		  rule {
			allowed_clients= "10.10.13.0"
			access= "ReadWrite"
			nfsv3 {
				checked =  true
			}
			nfsv4 {
				checked = false
			}
		  }
	    }
	  }
  `)
}

func testAccVolumeConfigUpdate() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "terraform-acceptance-test-1" {
		provider = netapp-gcp
		name = "terraform-acceptance-test-1"
		region = "us-west2"
		protocol_types = ["NFSv3"]
		network = "cvs-terraform-vpc"
		volume_path = "terraform-acceptance-test-path"
		size = 2048
		service_level = "extreme"
		snapshot_policy {
			enabled = true
			hourly_schedule {
			  snapshots_to_keep = 9
			  minute = 2
			}
			daily_schedule {
			  snapshots_to_keep = 20
			  hour = 22
			  minute = 10
			}
			weekly_schedule {
			  snapshots_to_keep = 15
			  hour = 2
			  minute = 35
			  day = "Tuesday"
			}
			monthly_schedule {
			  snapshots_to_keep = 10
			  hour = 3
			  minute = 5
			  days_of_month = 15
			}    
		}
		export_policy {
			rule {
			  allowed_clients = "10.0.0.0/8"
			  access= "ReadOnly"
			  nfsv3 {
				checked =  true
			  }
			  nfsv4 {
				checked = false
			  }
			}
		  rule {
			allowed_clients= "10.10.13.1"
			access= "ReadOnly"
			nfsv3 {
				checked =  true
			  }
			  nfsv4 {
				checked = false
			  }
			}
		  }
	  }
	`)
}

func testAccVolumeConfigCreateSMB() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "terraform-acceptance-test-1-SMB" {
		provider = netapp-gcp
		name = "terraform-acceptance-test-1-SMB"
		region = "us-east4"
		protocol_types = ["SMB"]
		network = "cvs-terraform-vpc"
		size = 1024
		service_level = "extreme"
		snapshot_policy {
		  enabled = true
		  daily_schedule {
			hour = 10
			minute = 1
		  }
		}
	  }
	`)
}

func testAccWaitSeconds(second int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(second) * time.Second)
		return nil
	}
}

func testCheckResourceAttr(name string, key string, value string) resource.TestCheckFunc {
	// use testCheckResourceAttr func in testing.go, add a 10 seconds sleep before returning a error.
	// Becuase the volume is still in transition between states thus not ready to delete yet.
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		log.Printf("rsrs rsrsrs rsrs: %#v", rs.Primary)
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}

		emptyCheck := false
		if value == "0" && (strings.HasSuffix(key, ".#") || strings.HasSuffix(key, ".%")) {
			emptyCheck = true
		}

		if v, ok := rs.Primary.Attributes[key]; !ok || v != value {
			if emptyCheck && !ok {
				return nil
			}

			time.Sleep(time.Duration(10) * time.Second)

			if !ok {
				return fmt.Errorf("%s: Attribute '%s' not found", name, key)
			}
			return fmt.Errorf(
				"%s: Attribute '%s' expected %#v, got %#v",
				name,
				key,
				value,
				v)
		}
		return nil
	}
}
