package gcp

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

// Provider is the main method for NetApp GCP Terraform provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GCP_PROJECT", nil),
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]+$"), "project must be an numerical project number"),
				Description:  "The project number for GCP API operations.",
			},
			"service_account": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GCP_SERVICE_ACCOUNT", nil),
				Description: "The private key path for GCP API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"netapp-gcp_volume":           resourceGCPVolume(),
			"netapp-gcp_active_directory": resourceGCPActiveDirectory(),
			"netapp-gcp_snapshot":         resourceGCPSnapshot(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"netapp-gcp_volume": dataSourceGCPVolume(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := configStuct{
		Project:        d.Get("project").(string),
		ServiceAccount: d.Get("service_account").(string),
	}

	return config.clientFun()
}
