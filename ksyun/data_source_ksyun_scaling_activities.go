package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunScalingActivities() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingActivitiesRead,
		Schema: map[string]*schema.Schema{
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

			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"scaling_activities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cause": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_activity_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"error_code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"success_instance_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingActivitiesRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunScalingActivities()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"scaling_group_id": {},
		"end_time":         {},
		"start_time":       {},
	}
	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingActivity list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingActivity", req)
		resp, err := client.kecconn.DescribeScalingActivity(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingActivity list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingActivitySet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, all)

	err = dataSourceKsyunScalingActivitiesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingActivity list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingActivitiesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingActivities()
	targetName := "scaling_activities"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingActivityId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, scalingActivitySpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}

func scalingActivitySpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["Desciption"] = SdkResponseMapping{Field: "description"}
	specialMapping["SuccInsList"] = SdkResponseMapping{Field: "success_instance_list"}
	return specialMapping
}
