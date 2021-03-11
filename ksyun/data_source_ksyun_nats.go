package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunNats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunNatsRead,
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

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"nats": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"nat_ip_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"band_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"nat_ip_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nat_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nat_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"associate_nat_set": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceKsyunNatsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var allNats []interface{}

	limit := 500
	offset := 1

	client := meta.(*KsyunClient)
	allNats = []interface{}{}
	result = []map[string]interface{}{}
	readNat := make(map[string]interface{})

	if ids, ok := d.GetOk("ids"); ok {
		SchemaSetToInstanceMap(ids, "NatId", &readNat)
	}

	if ids, ok := d.GetOk("project_ids"); ok {
		SchemaSetToInstanceMap(ids, "ProjectId", &readNat)
	}

	//filter
	index := 0
	if vpcIds, ok := d.GetOk("vpc_ids"); ok {
		index = index + 1
		SchemaSetToFilterMap(vpcIds, "vpc-id", index, &readNat)
	}

	for {
		readNat["MaxResults"] = limit
		readNat["NextToken"] = offset

		logger.Debug(logger.ReqFormat, "DescribeNats", readNat)
		resp, err := client.vpcconn.DescribeNats(&readNat)
		if err != nil {
			return fmt.Errorf("error on reading nat list req(%v):%v", readNat, err)
		}
		l := (*resp)["NatSet"].([]interface{})
		allNats = append(allNats, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range allNats {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["NatName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, allNats)
	}

	err := dataSourceKsyunNatsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading nat list, %s", err)
	}

	return nil
}

func dataSourceKsyunNatsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "NatId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item["NatId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			_, natIpSet, _ := SdkSliceMapping(nil, item["NatIpSet"].([]interface{}), SdkSliceData{
				SliceMappingFunc: func(next map[string]interface{}) map[string]interface{} {
					return map[string]interface{}{
						"nat_ip":    next["NatIp"],
						"nat_ip_id": next["NatIpId"],
					}
				},
			})

			if item["AssociateNatSet"] != nil {
				_, associateNatSet, _ := SdkSliceMapping(nil, item["AssociateNatSet"].([]interface{}), SdkSliceData{
					SliceMappingFunc: func(next map[string]interface{}) map[string]interface{} {
						return map[string]interface{}{
							"subnet_id": next["SubnetId"],
						}
					},
				})
				return map[string]interface{}{
					"id":                item["NatId"],
					"vpc_id":            item["VpcId"],
					"nat_name":          item["NatName"],
					"nat_mode":          item["NatMode"],
					"nat_type":          item["NatType"],
					"nat_ip_set":        natIpSet,
					"associate_nat_set": associateNatSet,
					"nat_ip_number":     item["NatIpNumber"],
					"band_width":        item["BandWidth"],
					"create_time":       item["CreateTime"],
				}
			}

			return map[string]interface{}{
				"id":            item["NatId"],
				"vpc_id":        item["VpcId"],
				"nat_name":      item["NatName"],
				"nat_mode":      item["NatMode"],
				"nat_type":      item["NatType"],
				"nat_ip_set":    natIpSet,
				"nat_ip_number": item["NatIpNumber"],
				"band_width":    item["BandWidth"],
				"create_time":   item["CreateTime"],
			}
		},
		TargetName: "nats",
	})
	return err
}
