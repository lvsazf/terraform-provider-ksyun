package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunEipsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"internet_gateway_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"instance_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"band_width_share_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"line_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"public_ip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ipv4",
					"ipv6",
					"all",
				}, false),
			},
			"eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allocation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"band_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"band_width_share_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_band_width_share": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunEipsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunEips()

	limit := 500
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform

	only = map[string]SdkReqTransform{
		"ids":                  {mapping: "AllocationId", Type: TransformWithN},
		"project_id":           {Type: TransformWithN},
		"network_interface_id": {Type: TransformWithFilter},
		"instance_type":        {Type: TransformWithFilter},
		"band_width_share_id":  {Type: TransformWithFilter},
		"line_id":              {Type: TransformWithFilter},
		"public_ip":            {Type: TransformWithFilter},
		"ip_version":           {},
	}

	req, err = SdkRequestAutoMappingNew(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading Addresses list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeAddresses", req)
		resp, err := client.eipconn.DescribeAddresses(&req)
		if err != nil {
			return fmt.Errorf("error on reading Addresses list req(%v):%v", req, err)
		}
		l := (*resp)["AddressesSet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, all)

	err = dataSourceKsyunEipsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading Addresses list, %s", err)
	}
	return nil
}

func dataSourceKsyunEipsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunEips()
	targetName := "eips"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "AllocationId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, nil, nil)
		},
		TargetName: targetName,
	})
	return err
}
