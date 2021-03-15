package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunScalingInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingInstancesRead,
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

			"scaling_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"health_status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"creation_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"scaling_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"add_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protected_from_detach": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingInstancesRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunScalingInstances()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"scaling_group_id": {},
		"health_status":    {},
		"creation_type":    {},
	}

	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingInstance list, %s", err)
	}

	if ids, ok := d.GetOk("scaling_instance_ids"); ok {
		SchemaSetToInstanceMap(ids, "ScalingInstanceId", &req)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingInstance", req)
		resp, err := client.kecconn.DescribeScalingInstance(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingInstance list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingInstanceSet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, all)

	err = dataSourceKsyunScalingInstancesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingInstance list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingInstancesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingInstances()
	targetName := "scaling_instances"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "InstanceId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + d.Get("scaling_group_id").(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, scalingInstanceSpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}

func scalingInstanceSpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["InstanceId"] = SdkResponseMapping{Field: "scaling_instance_id"}
	specialMapping["InstanceName"] = SdkResponseMapping{Field: "scaling_instance_name"}
	specialMapping["ProtectedFromScaleIn"] = SdkResponseMapping{Field: "protected_from_detach"}
	return specialMapping
}
