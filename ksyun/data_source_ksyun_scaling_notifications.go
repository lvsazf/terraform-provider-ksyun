package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunScalingNotifications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingNotificationsRead,
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

			"scaling_notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_notification_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_notification_types": {
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

func dataSourceKsyunScalingNotificationsRead(d *schema.ResourceData, meta interface{}) error {
	resource := dataSourceKsyunScalingNotifications()
	var result []map[string]interface{}
	var err error

	client := meta.(*KsyunClient)
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"scaling_group_id": {},
	}

	req, err = SdkRequestAutoMapping(d, resource, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingNotification list, %s", err)
	}

	logger.Debug(logger.ReqFormat, "DescribeScalingNotification", req)
	resp, err := client.kecconn.DescribeScalingNotification(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingNotification list req(%v):%v", req, err)
	}
	if (*resp)["ScalingNotificationSet"] == nil {
		return nil
	}
	l := (*resp)["ScalingNotificationSet"].([]interface{})

	merageResultDirect(&result, l)

	err = dataSourceKsyunScalingNotificationsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingNotification list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingNotificationsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingNotifications()
	targetName := "scaling_notifications"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingNotificationId",
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
