package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunScalingScheduledTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingScheduledTasksRead,
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

			"scaling_scheduled_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"scaling_scheduled_tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_scheduled_task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_scheduled_task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"readjust_max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"readjust_min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"readjust_expect_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"recurrence": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"repeat_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"repeat_cycle": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingScheduledTasksRead(d *schema.ResourceData, meta interface{}) error {
	resource := dataSourceKsyunScalingScheduledTasks()
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
		"ids":                         {mapping: "ScalingScheduledTaskId", Type: TransformWithN},
		"scaling_group_id":            {},
		"scaling_scheduled_task_name": {},
	}

	req, err = SdkRequestAutoMapping(d, resource, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingScheduledTask list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScheduledTask", req)
		resp, err := client.kecconn.DescribeScheduledTask(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingScheduledTask list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingScheduleTaskSet"].([]interface{})
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
			if r != nil && !r.MatchString(item["ScalingScheduledTaskName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, all)
	}

	err = dataSourceKsyunScalingScheduledTasksSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingScheduledTask list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingScheduledTasksSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingScheduledTasks()
	targetName := "scaling_scheduled_tasks"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingScheduledTaskId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, nil)
		},
		TargetName: targetName,
	})
	return err
}
