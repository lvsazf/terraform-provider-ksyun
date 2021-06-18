package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVpcsRead,

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
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunVpcsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var err error
	r := dataSourceKsyunVpcs()

	client := meta.(*KsyunClient)
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform

	only = map[string]SdkReqTransform{
		"ids": {mapping: "VpcId", Type: TransformWithN},
	}

	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading Addresses list, %s", err)
	}

	logger.Debug(logger.ReqFormat, "DescribeVpcs", req)
	resp, err := client.vpcconn.DescribeVpcs(&req)
	if err != nil {
		return fmt.Errorf("error on reading vpc list req(%v):%v", req, err)
	}
	l := (*resp)["VpcSet"].([]interface{})
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range l {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["VpcName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, l)
	}
	err = dataSourceKsyunVpcsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading vpc list, %s", err)
	}
	return nil
}

func dataSourceKsyunVpcsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunVpcs()
	targetName := "vpcs"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "VpcId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, resourceKsyunVpcExtra())
		},
		TargetName: targetName,
	})
	return err
}

func resourceKsyunVpcExtra() map[string]SdkResponseMapping {
	extra := make(map[string]SdkResponseMapping)
	extra["VpcName"] = SdkResponseMapping{
		Field:    "name",
		KeepAuto: true,
	}
	return extra
}
