package gcp

import (
	"context"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Provider is the main method for NetApp GCP Terraform provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GCP_PROJECT", nil),
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]+$|^[a-z][a-z0-9-]+[a-z0-9]$"),
					"Project format is not correct. It should be either numberical project number or project ID in xxx-xxx-xxx format."),
				Description: "The project number or project ID for GCP API operations.",
			},
			"service_account": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GCP_SERVICE_ACCOUNT", nil),
				Description: "The private key path for GCP API operations.",
			},
			"credentials": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GCP_CREDENTIALS", nil),
				Description: "The credentials for GCP API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"netapp-gcp_volume":             resourceGCPVolume(),
			"netapp-gcp_active_directory":   resourceGCPActiveDirectory(),
			"netapp-gcp_snapshot":           resourceGCPSnapshot(),
			"netapp-gcp_volume_backup":      resourceGCPVolumeBackup(),
			"netapp-gcp_volume_replication": resourceGCPVolumeReplication(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"netapp-gcp_volume":           dataSourceGCPVolume(),
			"netapp-gcp_active_directory": dataSourceGCPActiveDirectory(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func getProjectNumber(p string, d *schema.ResourceData) (string, error) {
	ctx := context.Background()
	var ts *google.Credentials
	var b []byte
	if v, ok := d.GetOk("service_account"); ok {
		var err error
		b, err = ioutil.ReadFile(v.(string))
		if err != nil {
			return "", err
		}
	} else if v, ok := d.GetOk("credentials"); ok {
		b = []byte(v.(string))
	}
	ts, err := google.CredentialsFromJSON(ctx, b, cloudresourcemanager.CloudPlatformScope)
	if err != nil {
		return "", err
	}
	c := oauth2.NewClient(ctx, ts.TokenSource)
	if err != nil {
		log.Printf("getProjectNumber: Not able to get client (%s)", err)
		return "", err
	}
	cloudresourcemanagerService, err := cloudresourcemanager.New(c)
	if err != nil {
		log.Printf("getProjectNumber: Cannot get cloud resource manager service(%s)", err)
		return "", err
	}
	resp, err := cloudresourcemanagerService.Projects.Get(p).Context(ctx).Do()
	if err != nil {
		log.Printf("getProjectNumber: Cannot find project number (%s)", err)
		return "", err
	}
	return strconv.FormatInt(resp.ProjectNumber, 10), nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var projectNumber string
	project := d.Get("project").(string)
	serviceAccount := d.Get("service_account").(string)
	credentials := d.Get("credentials").(string)
	// check if project is project number or project ID
	// project number is a string with numbers
	// project ID must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited
	isProjectNumber, _ := regexp.MatchString("^[0-9]+$", project)
	if !isProjectNumber {
		isProjectID, err := regexp.MatchString("^[a-z][a-z0-9-]+[a-z0-9]+$", project)
		if isProjectID {
			projectNumber, err = getProjectNumber(project, d)
			if err != nil {
				log.Printf("providerConfigure: Cannot find project number (%s)", err)
				return nil, err
			}
		} else {
			log.Printf("providerConfigure: Project %s format is not correct. It should be either numerical project number or project ID in xxx-xxx-xxxx format.", project)
			return nil, err
		}
	} else {
		projectNumber = project
	}
	config := configStuct{
		Project:        projectNumber,
		ServiceAccount: serviceAccount,
		Credentials:    credentials,
	}

	return config.clientFun()
}
