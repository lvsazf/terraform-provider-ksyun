package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunScalingGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scaling_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"remove_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"security_group_id_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"desired_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"subnet_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"slb_config_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"slb_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"listener_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"server_port_set": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
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

func dataSourceKsyunScalingGroupsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunScalingGroups()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"ids":                      {mapping: "ScalingGroupId", Type: TransformWithN},
		"scaling_group_name":       {},
		"scaling_configuration_id": {},
		"vpc_id":                   {},
	}

	req, err = SdkRequestAutoMapping(d, r, false, only, resourceKsyunScalingGroupExtra(d, false))
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroup list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingGroup", req)
		resp, err := client.kecconn.DescribeScalingGroup(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingGroup list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingGroupSet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range all {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["ScalingGroupName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, all)
	}

	err = dataSourceKsyunScalingGroupsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroup list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingGroupsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingGroups()
	targetName := "scaling_groups"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingGroupId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, nil)
		},
		TargetName: targetName,
	})
	return err
}
