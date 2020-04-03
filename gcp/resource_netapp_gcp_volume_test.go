package gcp

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVolume_basic(t *testing.T) {

	var volume listVolumeResult
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
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "service_level", "premium"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "protocol_types.0", "NFSv3"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.hour", "10"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.minute", "1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.snapshots_to_keep", "0"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.access", "ReadWrite"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.allowed_clients", "0.0.0.0/0"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.nfsv3.0.checked", "true"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.nfsv4.0.checked", "false"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.access", "ReadWrite"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.allowed_clients", "10.10.13.0"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.nfsv3.0.checked", "true"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.nfsv4.0.checked", "false"),
				),
			},
			{
				Config: testAccVolumeConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPVolumeExists("netapp-gcp_volume.terraform-acceptance-test-1", &volume),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "name", "terraform-acceptance-test-1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "size", "2048"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "region", "us-west2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "service_level", "standard"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.hour", "20"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.minute", "30"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.snapshots_to_keep", "0"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.access", "ReadOnly"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.allowed_clients", "10.0.0.0/8"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.nfsv3.0.checked", "false"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.nfsv4.0.checked", "true"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.access", "ReadOnly"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.allowed_clients", "10.10.13.1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.nfsv3.0.checked", "false"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.nfsv4.0.checked", "true"),
				),
			},
			{
				Config: testAccVolumeConfigCreateSMB(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPVolumeExists("netapp-gcp_volume.terraform-acceptance-test-1", &volume),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "name", "terraform-acceptance-test-1-SMB"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "size", "1024"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "region", "us-west2"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "service_level", "premium"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "protocol_types.0", "SMB"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.hour", "10"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.minute", "1"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "snapshot_policy.0.daily_schedule.0.snapshots_to_keep", "0"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.access", "ReadWrite"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.0.allowed_clients", "0.0.0.0/0"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.access", "ReadWrite"),
					testCheckResourceAttr("netapp-gcp_volume.terraform-acceptance-test-1", "export_policy.0.rule.1.allowed_clients", "10.10.13.0"),
				),
			},
		},
	})
}

func testAccCheckGCPVolumeDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "netapp-gcp_volume" {
			continue
		}
		response, err := client.getVolumeByID(listVolumesRequest{
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

func testAccCheckGCPVolumeExists(name string, volume *listVolumeResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}
		response, err := client.getVolumeByID(listVolumesRequest{
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
		service_level = "premium"
		snapshot_policy {
		  enabled = true
		  daily_schedule {
			hour = 10	
			minute = 1
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
		size = 2048
		service_level = "standard"
		snapshot_policy {
		  enabled = true
		  daily_schedule {
			hour = 20	
			minute = 30
		  }
		}
		export_policy {
			rule {
			  allowed_clients = "10.0.0.0/8"
			  access= "ReadOnly"
			  nfsv3 {
				checked =  false
			  }
			  nfsv4 {
				checked = true
			  }
			}
		  rule {
			allowed_clients= "10.10.13.1"
			access= "ReadOnly"
			nfsv3 {
				checked =  false
			  }
			  nfsv4 {
				checked = true
			  }
			}
		  }
	  }
	`)
}

func testAccVolumeConfigCreateSMB() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_volume" "terraform-acceptance-test-1" {
		provider = netapp-gcp
		name = "terraform-acceptance-test-1-SMB"
		region = "us-west2"
		protocol_types = ["SMB"]
		network = "cvs-terraform-vpc"
		size = 1024
		service_level = "premium"
		snapshot_policy {
		  enabled = true
		  daily_schedule {
			hour = 10	
			minute = 1
		  }
		}
		export_policy {
			rule {
			  allowed_clients = "0.0.0.0/0"
			  access= "ReadWrite"
			}
		  rule {
			allowed_clients= "10.10.13.0"
			access= "ReadWrite"
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
