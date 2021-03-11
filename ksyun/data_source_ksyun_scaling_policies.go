package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunScalingPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingPoliciesRead,
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

			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"scaling_policies_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"scaling_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"adjustment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"adjustment_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cool_down": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"dimension_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"comparison_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"repeat_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"function": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	resource := dataSourceKsyunScalingPolicies()
	var result []map[string]interface{}
	var all []interface{}
	var err error

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"ids":                   {mapping: "ScalingPolicyId", Type: TransformWithN},
		"scaling_group_id":      {},
		"scaling_policies_name": {},
	}

	req, err = SdkRequestAutoMapping(d, resource, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingPolicy list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingPolicy", req)
		resp, err := client.kecconn.DescribeScalingPolicy(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingPolicy list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingPolicySet"].([]interface{})
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
			if r != nil && !r.MatchString(item["ScalingPolicyName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, all)
	}

	err = dataSourceKsyunScalingPoliciesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingPolicy list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingPoliciesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingPolicies()
	targetName := "scaling_policies"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingPolicyId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			var compute map[string]interface{}
			if item["Metric"] != nil {
				compute, _ = SdkMapMapping(item["Metric"].(map[string]interface{}), SdkSliceData{
					SliceMappingFunc: func(m map[string]interface{}) map[string]interface{} {
						return SdkResponseAutoMapping(resource, targetName, m, nil, nil, nil)
					},
				})
			}
			return SdkResponseAutoMapping(resource, targetName, item, compute, nil, nil)
		},
		TargetName: targetName,
	})
	return err
}
