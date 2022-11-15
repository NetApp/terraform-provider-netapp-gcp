package gcp

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGCPVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGCPVolumeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_dp": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_vpc_project_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"snapshot_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"daily_schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"minute": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"hourly_schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minute": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"monthly_schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_of_month": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hour": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"minute": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"weekly_schedule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hour": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"minute": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"export_policy": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allowed_clients": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"has_root_access": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kerberos5_readonly": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"kerberos5_readwrite": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"kerberos5i_readonly": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"kerberos5i_readwrite": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"kerberos5p_readonly": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"kerberos5p_readwrite": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nfsv3": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"checked": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"nfsv4": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"checked": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"delete_on_creation_error": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regional_ha": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"snapshot_directory": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smb_share_settings": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"unix_permissions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_style": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_label": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGCPVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume: %#v", d)
	client := meta.(*Client)

	volume := volumeRequest{}

	volume.Name = d.Get("name").(string)
	volume.Region = d.Get("region").(string)

	// Resolve volume name to volume UUID
	var res volumeResult
	res, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		return err
	}

	// Set ID to volume UUID and use normal volume read call to do parsing of attributes
	d.SetId(res.VolumeID)
	return resourceGCPVolumeRead(d, meta)
}
