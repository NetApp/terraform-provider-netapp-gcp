package gcp

import (
	"fmt"
	"log"
	"strings"

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
										Type:     schema.TypeBool,
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
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
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

	var res volumeResult
	res, err := client.getVolumeByNameOrCreationToken(volume)
	if err != nil {
		return err
	}

	d.SetId(res.VolumeID)
	if err := d.Set("name", res.Name); err != nil {
		return fmt.Errorf("Error reading volume name: %s", err)
	}
	if err := d.Set("type_dp", res.TypeDP); err != nil {
		return fmt.Errorf("Error reading type_dp: %s", err)
	}
	if err := d.Set("size", res.Size/GiBToBytes); err != nil {
		return fmt.Errorf("Error reading volume size: %s", err)
	}
	if err := d.Set("service_level", res.ServiceLevel); err != nil {
		return fmt.Errorf("Error reading volume service_level: %s", err)
	}
	for i, protocol := range res.ProtocolTypes {
		if protocol == "CIFS" {
			res.ProtocolTypes[i] = "SMB"
		}
	}
	if err := d.Set("protocol_types", res.ProtocolTypes); err != nil {
		return fmt.Errorf("Error reading volume protocol_types: %s", err)
	}
	if err := d.Set("volume_path", res.CreationToken); err != nil {
		return fmt.Errorf("Error reading volume path or Creation Token: %s", err)
	}
	network := res.Network
	index := strings.Index(network, "networks/")
	if index > -1 {
		network = network[index+len("networks/"):]
	}
	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("Error reading volume network: %s", err)
	}
	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("Error reading volume region: %s", err)
	}
	snapshotPolicy := flattenSnapshotPolicy(res.SnapshotPolicy)
	exportPolicy := flattenExportPolicy(res.ExportPolicy)
	if err := d.Set("snapshot_policy", snapshotPolicy); err != nil {
		return fmt.Errorf("Error reading volume snapshot_policy: %s", err)
	}
	if len(res.ExportPolicy.Rules) > 0 {
		if err := d.Set("export_policy", exportPolicy); err != nil {
			return fmt.Errorf("Error reading volume export_policy: %s", err)
		}
	} else {
		a := schema.NewSet(schema.HashString, []interface{}{})
		if err := d.Set("export_policy", a); err != nil {
			return fmt.Errorf("Error reading volume export_policy: %s", err)
		}
	}
	mountPoints := flattenMountPoints(res.MountPoints)
	if err := d.Set("mount_points", mountPoints); err != nil {
		return fmt.Errorf("Error reading volume mount_points: %s", err)
	}
	if err := d.Set("zone", res.Zone); err != nil {
		return fmt.Errorf("Error reading zone: %s", err)
	}
	if err := d.Set("storage_class", res.StorageClass); err != nil {
		return fmt.Errorf("Error reading storage_class: %s", err)
	}
	return nil
}
