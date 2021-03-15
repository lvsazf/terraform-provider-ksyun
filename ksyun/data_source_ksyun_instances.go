package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"search": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc: validation.ValidateRegexp,
			},

			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_interface": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"network_interface_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"group_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"instance_state": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_configure": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"v_c_p_u": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"g_p_u": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"memory_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"data_disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"data_disk_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"root_disk_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_state": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitoring": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"sriov_net_support": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_interface_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_interface_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_set": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"security_group_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"group_set": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"d_n_s1": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"d_n_s2": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"system_disk": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stopped_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_what": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auto_scaling_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_show_sriov_net_support": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"data_disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"delete_with_instance": {
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
	}
}

func dataSourceKsyunInstancesRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunInstances()

	limit := 100
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform

	only = map[string]SdkReqTransform{
		"ids":               {mapping: "InstanceId", Type: TransformWithN},
		"project_id":        {Type: TransformWithN},
		"vpc_id":            {Type: TransformWithFilter},
		"subnet_id":         {Type: TransformWithFilter},
		"search":            {mapping: "Search"},
		"network_interface": {Type: TransformListFilter},
		"instance_state":    {Type: TransformListFilter},
		"availability_zone": {mappings: map[string]string{
			"availability_zone.name": "availability-zone-name",
		}, Type: TransformListFilter},
	}

	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading Instance list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeInstances", req)
		resp, err := client.kecconn.DescribeInstances(&req)
		if err != nil {
			return fmt.Errorf("error on reading Instance list req(%v):%v", req, err)
		}
		l := (*resp)["InstancesSet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, all)

	err = dataSourceKsyunInstancesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading Instance list, %s", err)
	}
	return nil
}

func dataSourceKsyunInstancesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunInstances()
	targetName := "instances"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "InstanceId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, kecInstanceSpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}

func kecInstanceSpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["KeySet"] = SdkResponseMapping{Field: "key_id"}
	return specialMapping
}
