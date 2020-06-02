package gcp

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccActiveDirectory_basic(t *testing.T) {
	// if the active direactory already exists, the acceptance test will fail.

	var activeDirectory listActiveDirectoryResult

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCPActiveDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryConfigCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPActiveDirectoryExists("netapp-gcp_active_directory.terraform-acceptance-test-1", &activeDirectory),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "username", "test_user"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "region", "us-central1"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "domain", "example.com"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "dns_server", "10.0.0.0"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "net_bios", "cvserver"),
				),
			},
			{
				Config: testAccActiveDirectoryConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPActiveDirectoryExists("netapp-gcp_active_directory.terraform-acceptance-test-1", &activeDirectory),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "username", "new_test_user"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "region", "us-central1"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "domain", "newExample.com"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "dns_server", "10.0.0.1"),
					resource.TestCheckResourceAttr("netapp-gcp_active_directory.terraform-acceptance-test-1", "net_bios", "cvservers"),
					testAccWaitSeconds(10),
				),
			},
		},
	})
}

func testAccCheckGCPActiveDirectoryDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "netapp-gcp_active_directory" {
			continue
		}
		response, err := client.listActiveDirectoryForRegion(listActiveDirectoryRequest{
			UUID:   rs.Primary.ID,
			Region: rs.Primary.Attributes["region"],
		})
		if err == nil {
			if response.UUID != "" {
				return fmt.Errorf("Active directory (%s) still exists.", response.UUID)
			}

		}
	}

	return nil
}

func testAccCheckGCPActiveDirectoryExists(name string, activeDirectory *listActiveDirectoryResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No active directory ID is set")
		}
		response, err := client.listActiveDirectoryForRegion(listActiveDirectoryRequest{
			UUID:   rs.Primary.ID,
			Region: rs.Primary.Attributes["region"],
		})

		if err != nil {
			return err
		}

		if response.UUID != rs.Primary.ID {
			return fmt.Errorf("Resource ID and active directory ID do not match")
		}

		*activeDirectory = response

		return nil
	}
}

func testAccActiveDirectoryConfigCreate() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_active_directory" "terraform-acceptance-test-1" {
		provider = netapp-gcp
		region = "us-central1"
		  username = "test_user"
		  password = "netapp"
		domain = "example.com"
		dns_server = "10.0.0.0"
		net_bios = "cvserver"
	  }
	`)
}

func testAccActiveDirectoryConfigUpdate() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_active_directory" "terraform-acceptance-test-1" {
		provider = netapp-gcp
		region = "us-central1"
		  username = "new_test_user"
		  password = "netapp"
		domain = "newExample.com"
		dns_server = "10.0.0.1"
		net_bios = "cvservers"
	  }
	`)
}
