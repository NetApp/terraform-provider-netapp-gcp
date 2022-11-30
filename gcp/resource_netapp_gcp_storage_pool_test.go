package gcp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Default go test timeout is 10min, which is too short for this test. Set to 20min
// TODO: deleteStoragePool doesn't wait for completion. Therefore, CheckDestroy doesn't make sense yet
// TODO: Fix issue #79 to be able to add additional test for regional pool
// TODO: Add test for regional pool

func TestAccStoragePoolZonal(t *testing.T) {
	var pool storagePool

	// Uncomment to inherit global constants
	// region := Region
	// network := Network
	// Overwrites for my environment
	region := "europe-west1"
	network := "cvs-prd-shared-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckGCPStoragePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStoragePoolConfigZonalCreate(region, network),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGCPStoragePoolExists("netapp-gcp_storage_pool.storage-pool-zonal", &pool),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "name", "storage-pool-zonal"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "size", "1024"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "region", region),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "zone", region+"-b"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "network", network),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "service_level", "StandardSW"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "storage_class", "software"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "global_ad_access", "false"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "managed_pool", "false"),
					// resource.TestCheckTypeSetElemNestedAttrs(
					// 	"netapp-gcp_storage_pool.storage-pool-zonal",
					// 	"billing_label.0",
					// 	map[string]string{
					// 		"key":   "department",
					// 		"value": "1234",
					// 	},
					// ),
				),
			},
			{
				Config: testAccStoragePoolConfigZonalUpdate(region, network),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGCPStoragePoolExists("netapp-gcp_storage_pool.storage-pool-zonal", &pool),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "name", "storage-pool-zonal"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "size", "5000"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "region", region),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "zone", region+"-b"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "network", network),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "service_level", "StandardSW"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "storage_class", "software"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "global_ad_access", "true"),
					resource.TestCheckResourceAttr("netapp-gcp_storage_pool.storage-pool-zonal", "managed_pool", "false"),
					// resource.TestCheckTypeSetElemNestedAttrs(
					// 	"netapp-gcp_storage_pool.storage-pool-zonal",
					// 	"billing_label.0",
					// 	map[string]string{
					// 		"key":   "department",
					// 		"value": "5678",
					// 	},
					// ),
				),
			},
		},
	})
}

func testAccCheckGCPStoragePoolExists(name string, pool *storagePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pool ID is set")
		}
		response, err := client.getStoragePoolByID(&storagePool{
			PoolID: rs.Primary.ID,
			Region: rs.Primary.Attributes["region"],
		})

		if err != nil {
			return err
		}

		if response.PoolID != rs.Primary.ID {
			return fmt.Errorf("Resource ID and pool ID do not match")
		}

		*pool = response

		return nil
	}
}

func testAccStoragePoolConfigZonalCreate(region, network string) string {
	return fmt.Sprintf(`
	resource "netapp-gcp_storage_pool" "storage-pool-zonal" {
		name = "storage-pool-zonal"
		region = "%s"
		zone = "%s"
		network = "%s"
		size = 1024
		service_level = "StandardSW"
		storage_class = "software"
		global_ad_access = false
		billing_label {
		  key = "department"
		  value = "1234"
		}
	  }
  `, region, region+"-b", network)
}

func testAccStoragePoolConfigZonalUpdate(region, network string) string {
	return fmt.Sprintf(`
	resource "netapp-gcp_storage_pool" "storage-pool-zonal" {
		name = "storage-pool-zonal"
		region = "%s"
		zone = "%s"
		network = "%s"
		size = 5000
		service_level = "StandardSW"
		storage_class = "software"
		global_ad_access = true
		billing_label {
		  key = "department"
		  value = "5678"
		}
	  }
  `, region, region+"-b", network)
}
