package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/netapp/terraform-provider-netapp-gcp/gcp"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gcp.Provider,
	})
}
