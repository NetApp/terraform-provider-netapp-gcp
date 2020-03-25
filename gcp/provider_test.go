package gcp

import (
	"testing"

	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"netapp-gcp": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GCP_PROJECT"); v == "" {
		t.Fatal("GCP_PROJECT must be set for acceptance tests")
	}

	if v := os.Getenv("GCP_SERVICE_ACCOUNT"); v == "" {
		t.Fatal("GCP_SERVICE_ACCOUNT must be set for acceptance tests")
	}

}
