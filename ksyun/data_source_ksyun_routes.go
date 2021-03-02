package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunRoutesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
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

			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"routes": {
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

						"destination_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"next_hop_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"gateway_name": {
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

func dataSourceKsyunRoutesRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var allRoutes []interface{}

	limit := 500
	offset := 1

	client := meta.(*KsyunClient)
	allRoutes = []interface{}{}
	result = []map[string]interface{}{}
	readRoute := make(map[string]interface{})

	if ids, ok := d.GetOk("ids"); ok {
		SchemaSetToInstanceMap(ids, "RouteId", &readRoute)
	}

	//filter
	index := 0
	if vpcIds, ok := d.GetOk("vpc_ids"); ok {
		index = index + 1
		SchemaSetToFilterMap(vpcIds, "vpc-id", index, &readRoute)
	}

	if vpcIds, ok := d.GetOk("instance_ids"); ok {
		index = index + 1
		SchemaSetToFilterMap(vpcIds, "instance-id", index, &readRoute)
	}

	for {
		readRoute["MaxResults"] = limit
		readRoute["NextToken"] = offset

		logger.Debug(logger.ReqFormat, "DescribeRoutes", readRoute)
		resp, err := client.vpcconn.DescribeRoutes(&readRoute)
		if err != nil {
			return fmt.Errorf("error on reading route list req(%v):%v", readRoute, err)
		}
		l := (*resp)["RouteSet"].([]interface{})
		allRoutes = append(allRoutes, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, allRoutes)

	err := dataSourceKsyunRoutesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading route list, %s", err)
	}

	return nil
}

func dataSourceKsyunRoutesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "RouteId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item["RouteId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			_, aaa, _ := SdkSliceMapping(nil, item["NextHopSet"].([]interface{}), SdkSliceData{
				SliceMappingFunc: func(next map[string]interface{}) map[string]interface{} {
					return map[string]interface{}{
						"gateway_id":   next["GatewayId"],
						"gateway_name": next["GatewayName"],
					}
				},
			})

			return map[string]interface{}{
				"id":                     item["RouteId"],
				"vpc_id":                 item["VpcId"],
				"destination_cidr_block": item["DestinationCidrBlock"],
				"route_type":             item["RouteType"],
				"next_hop_set":           aaa,
				"create_time":            item["CreateTime"],
			}
		},
		TargetName: "routes",
	})
	return err
}
