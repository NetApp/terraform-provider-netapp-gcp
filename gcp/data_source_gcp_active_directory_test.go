package gcp

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceActiveDirectory_basic(t *testing.T) {
	datasourceName := "data.netapp-gcp_active_directory.ad-us-west2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// if create a new active directory before creating the data source, the acceptance test passes locally but fails on Jekins server.
			// Currently, data source is created based on existing active directory.
			// {
			// 	Config: testAccActiveDirectorResource(),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckGCPActiveDirectoryExists("netapp-gcp_active_directory.terraform-acceptance-test-1", &activeDirectory),
			// 	),
			// },
			{
				Config: testAccActiveDirectoryDataResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "username", "bsuhas"),
					resource.TestCheckResourceAttr(datasourceName, "dns_server", "10.168.0.6"),
					resource.TestCheckResourceAttr(datasourceName, "netbios", "cvs-smb"),
				),
			},
		},
	})
}

func testAccActiveDirectorResource() string {
	return fmt.Sprintf(`
	resource "netapp-gcp_active_directory" "terraform-acceptance-test-1" {
		provider = netapp-gcp
		region = "us-central1"
		  username = "test_user"
		  password = "netapp"
		domain = "example.com"
		dns_server = "10.0.0.0"
		net_bios = "cvserver"
		organizational_unit = "CN=Computers"
	  }

	`)
}

func testAccActiveDirectoryDataResource() string {
	return fmt.Sprintf(`
	data "netapp-gcp_active_directory" "ad-us-west2" {
		region = "us-west2"
	}
	`)

}
